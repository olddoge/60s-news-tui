package app

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/ui"
)

// View renders the current page.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	switch m.page {
	case PageLoading:
		return m.viewLoading()
	case PageEndpointList:
		return m.viewEndpointList()
	case PageEncodingSelect:
		return m.viewEncodingSelect()
	case PageResult:
		return m.viewResult()
	case PageSettings:
		return m.viewSettings()
	case PageError:
		return m.viewError()
	default:
		return ""
	}
}

func (m Model) viewLoading() string {
	message := m.loadingMessage
	if message == "" {
		message = "Loading..."
	}
	return ui.ContainerStyle.Render(
		ui.LoadingStyle.Render(message),
	)
}

func (m Model) viewEndpointList() string {
	var b strings.Builder

	b.WriteString(ui.RenderTitle("Endpoint TUI"))
	b.WriteString("\n")

	baseURL := m.config.BaseURL
	if baseURL == "" {
		baseURL = "(not configured)"
	}
	b.WriteString(ui.InfoStyle.Render("Base URL: " + baseURL))
	b.WriteString("\n\n")

	b.WriteString(ui.InfoStyle.Render("Endpoints:"))
	b.WriteString("\n\n")

	listHeight := m.height - 10
	if listHeight < 1 {
		listHeight = 1
	}

	start := 0
	if m.cursor >= listHeight {
		start = m.cursor - listHeight + 1
	}

	for i := start; i < len(m.endpoints) && i < start+listHeight; i++ {
		ep := m.endpoints[i]
		if i == m.cursor {
			b.WriteString(ui.SelectedStyle.Render("> " + formatEndpointLine(ep, m.width-4)))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + formatEndpointLine(ep, m.width-4)))
		}
		b.WriteString("\n")
	}

	if len(m.endpoints) > listHeight {
		b.WriteString(ui.HelpStyle.Render(
			fmt.Sprintf("  ... %d/%d", m.cursor+1, len(m.endpoints)),
		))
		b.WriteString("\n")
	}

	b.WriteString(ui.RenderHelp([]string{
		"up/down select",
		"Enter request",
		"s settings",
		"r refresh",
		"q quit",
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewEncodingSelect() string {
	var b strings.Builder

	ep := m.SelectedEndpoint()
	epPath := ""
	if ep != nil {
		epPath = ep.Path
	}

	b.WriteString(ui.RenderTitle("Request Options"))
	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render("Endpoint: " + epPath))
	b.WriteString("\n\n")
	b.WriteString(ui.InfoStyle.Render("Select response format:"))
	b.WriteString("\n\n")

	for i, enc := range m.encodings {
		if i == m.encodingCursor {
			b.WriteString(ui.SelectedStyle.Render("> " + enc))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + enc))
		}
		b.WriteString("\n")
	}

	b.WriteString(ui.RenderHelp([]string{
		"up/down select",
		"Enter run",
		"Esc back",
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewResult() string {
	var b strings.Builder

	r := m.result
	ep := m.SelectedEndpoint()
	epPath := ""
	if ep != nil {
		epPath = ep.Name
	}

	b.WriteString(ui.RenderTitle("Request Complete"))
	b.WriteString("\n")

	infoLines := []string{
		ui.LabelStyle.Render("Endpoint:") + ui.ValueStyle.Render(epPath),
		ui.LabelStyle.Render("Format:") + ui.ValueStyle.Render(m.SelectedEncoding()),
		ui.LabelStyle.Render("URL:") + ui.ValueStyle.Render(ui.Truncate(r.URL, m.width-10)),
	}

	if r.Cancelled {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("Status:")+ui.WarningStyle.Render("cancelled"),
		)
	} else if r.Error != nil || r.ExitCode != 0 {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("Status:")+ui.ErrorStyle.Render(fmt.Sprintf("failed (exit=%d)", r.ExitCode)),
		)
		if r.Stderr != "" {
			infoLines = append(infoLines,
				ui.LabelStyle.Render("Error:")+ui.ErrorStyle.Render(ui.Truncate(r.Stderr, m.width-10)),
			)
		}
	} else {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("Duration:")+ui.ValueStyle.Render(r.Duration.Truncate(0).String()),
			ui.LabelStyle.Render("Status:")+ui.SuccessStyle.Render("success"),
		)
	}

	for _, line := range infoLines {
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render("Response:"))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())

	b.WriteString(ui.RenderHelp([]string{
		"up/down/PgUp/PgDn scroll",
		"Home/End jump",
		"r retry",
		"b/Esc back",
		"q quit",
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewSettings() string {
	var b strings.Builder

	b.WriteString(ui.RenderTitle("Settings"))
	if m.config.BaseURL == "" {
		b.WriteString("\n")
		b.WriteString(ui.WarningStyle.Render("Configure a base URL before using Endpoint TUI. Example: http://127.0.0.1:8080"))
	}
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render("Base URL:"))
	b.WriteString("\n")
	b.WriteString(m.settingsBaseURL.View())
	b.WriteString("\n\n")

	b.WriteString(ui.LabelStyle.Render("Default:"))
	b.WriteString("\n")
	for i, enc := range m.encodings {
		if i == m.settingsEncodingCursor {
			b.WriteString(ui.SelectedStyle.Render("> " + enc))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + enc))
		}
		b.WriteString("\n")
	}

	if m.settingsValidationError != "" {
		b.WriteString("\n")
		b.WriteString(ui.ErrorStyle.Render(m.settingsValidationError))
	}

	if m.settingsSaved {
		b.WriteString("\n")
		b.WriteString(ui.SuccessStyle.Render("Config saved"))
	}

	b.WriteString(ui.RenderHelp([]string{
		"up/down select format",
		"Ctrl+S save",
		"Esc cancel",
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(ui.RenderTitle("Endpoint TUI"))
	b.WriteString("\n")

	if m.loadErr != nil {
		b.WriteString(ui.ErrorStyle.Render("Load failed: "))
		b.WriteString(ui.ValueStyle.Render(safeErrorMessage(m.loadErr)))
	}

	b.WriteString(ui.RenderHelp([]string{
		"r reload",
		"s settings",
		"q quit",
	}))

	return ui.ContainerStyle.Render(b.String())
}

func formatEndpointLine(ep api.Endpoint, width int) string {
	path := ep.Path
	name := ep.Name

	if name == path {
		return ui.Truncate(path, width)
	}

	line := ui.PadRight(name, 20) + " " + path
	return ui.Truncate(line, width)
}

func safeErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	msg = regexp.MustCompile(`https?://[^\s"]+`).ReplaceAllString(msg, "[endpoint service]")
	msg = regexp.MustCompile(`\[[^\]]+\]:\d+`).ReplaceAllString(msg, "[address]")
	msg = regexp.MustCompile(`\b([A-Za-z0-9.-]+):\d+\b`).ReplaceAllString(msg, "$1")
	return msg
}

func formatResultContent(r api.CurlResult, encoding string) string {
	if r.Cancelled {
		return "request cancelled"
	}

	if r.Error != nil && r.ExitCode != 0 {
		content := r.Error.Error()
		if r.Stderr != "" {
			content += "\n\nstderr:\n" + r.Stderr
		}
		if r.Stdout != "" {
			content += "\n\nstdout:\n" + r.Stdout
		}
		return content
	}

	content := r.Stdout
	if content == "" {
		return "(empty response)"
	}

	switch encoding {
	case "json":
		return formatJSON(content)
	case "markdown", "text":
		return sanitizeContent(content)
	default:
		return sanitizeContent(content)
	}
}

func formatJSON(content string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(content), &obj); err != nil {
		return sanitizeContent(content)
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return sanitizeContent(content)
	}
	return string(formatted)
}

func sanitizeContent(content string) string {
	var b strings.Builder
	b.Grow(len(content))

	for _, r := range content {
		if r == '\n' || r == '\t' || r == '\r' {
			b.WriteRune(r)
			continue
		}
		if r < 32 || (r >= 0x7F && r <= 0x9F) {
			continue
		}
		if !utf8.ValidRune(r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
