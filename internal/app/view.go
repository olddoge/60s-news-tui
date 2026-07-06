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
	var b strings.Builder
	message := m.loadingMessage
	if message == "" {
		message = "Loading..."
	}
	b.WriteString(m.brandHeader())
	b.WriteString(ui.LoadingStyle.Render(message))
	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewEndpointList() string {
	var b strings.Builder

	b.WriteString(m.brandHeader())

	baseURL := m.config.BaseURL
	if baseURL == "" {
		baseURL = m.text("(not configured)", "（未配置）")
	}
	b.WriteString(ui.InfoStyle.Render(m.text("Endpoint: ", "接口地址：") + baseURL))
	b.WriteString("\n\n")

	b.WriteString(ui.InfoStyle.Render(m.text("Endpoints:", "接口列表：")))
	b.WriteString("\n\n")

	if m.searching || m.search != "" {
		prefix := m.text("Search: ", "搜索：")
		if m.searching {
			prefix = "/ " + prefix
		}
		b.WriteString(ui.InfoStyle.Render(prefix + m.search))
		b.WriteString("\n\n")
	}

	listHeight := m.height - 13
	if listHeight < 1 {
		listHeight = 1
	}

	start := 0
	if m.cursor >= listHeight {
		start = m.cursor - listHeight + 1
	}

	indexes := m.filteredEndpointIndexes()
	if len(indexes) == 0 {
		b.WriteString(ui.WarningStyle.Render(m.text("No matching endpoints", "未找到匹配接口")))
		b.WriteString("\n")
	}

	for i := start; i < len(indexes) && i < start+listHeight; i++ {
		ep := m.endpoints[indexes[i]]
		if i == m.cursor {
			b.WriteString(ui.SelectedStyle.Render("> " + formatEndpointLine(i+1, ep, m.width-4)))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + formatEndpointLine(i+1, ep, m.width-4)))
		}
		b.WriteString("\n")
	}

	if len(indexes) > listHeight {
		b.WriteString(ui.HelpStyle.Render(
			fmt.Sprintf("  ... %d/%d", m.cursor+1, len(indexes)),
		))
		b.WriteString("\n")
	}

	if m.searching {
		b.WriteString(ui.RenderHelp([]string{
			m.text("⌨️ type search", "⌨️ 输入搜索"),
			m.text("✅ Enter finish", "✅ Enter 完成"),
			m.text("🧹 Esc clear", "🧹 Esc 清空"),
		}))
	} else {
		b.WriteString(ui.RenderHelp([]string{
			m.text("⬆️⬇️ up/down select", "⬆️⬇️ 上/下 选择"),
			m.text("🔎 / search", "🔎 / 搜索"),
			m.text("🚀 Enter request", "🚀 Enter 请求"),
			m.text("⚙️ s settings", "⚙️ s 设置"),
			m.text("🔄 r refresh", "🔄 r 刷新"),
			m.text("🚪 q quit", "🚪 q 退出"),
		}))
	}

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewEncodingSelect() string {
	var b strings.Builder

	ep := m.SelectedEndpoint()
	epPath := ""
	if ep != nil {
		epPath = ep.Path
	}

	b.WriteString(m.brandHeader())
	b.WriteString(ui.RenderTitle(m.text("Request Options", "请求选项")))
	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render(m.text("Endpoint: ", "接口：") + epPath))
	b.WriteString("\n\n")
	b.WriteString(ui.InfoStyle.Render(m.text("Select response format:", "请选择返回格式：")))
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
		m.text("⬆️⬇️ up/down select", "⬆️⬇️ 上/下 选择"),
		m.text("🚀 Enter run", "🚀 Enter 执行"),
		m.text("↩️ Esc back", "↩️ Esc 返回"),
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

	b.WriteString(m.brandHeader())
	b.WriteString(ui.RenderTitle(m.text("Request Complete", "请求完成")))
	b.WriteString("\n")

	infoLines := []string{
		ui.LabelStyle.Render(m.text("Endpoint:", "接口：")) + ui.ValueStyle.Render(epPath),
		ui.LabelStyle.Render(m.text("Format:", "格式：")) + ui.ValueStyle.Render(m.SelectedEncoding()),
		ui.LabelStyle.Render(m.text("URL:", "地址：")) + ui.ValueStyle.Render(ui.Truncate(r.URL, m.width-10)),
	}

	if r.Cancelled {
		infoLines = append(infoLines,
			ui.LabelStyle.Render(m.text("Status:", "状态："))+ui.WarningStyle.Render(m.text("cancelled", "已取消")),
		)
	} else if r.Error != nil || r.ExitCode != 0 {
		infoLines = append(infoLines,
			ui.LabelStyle.Render(m.text("Status:", "状态："))+ui.ErrorStyle.Render(fmt.Sprintf(m.text("failed (exit=%d)", "失败（退出码=%d）"), r.ExitCode)),
		)
		if r.Stderr != "" {
			infoLines = append(infoLines,
				ui.LabelStyle.Render(m.text("Error:", "错误："))+ui.ErrorStyle.Render(ui.Truncate(r.Stderr, m.width-10)),
			)
		}
	} else {
		infoLines = append(infoLines,
			ui.LabelStyle.Render(m.text("Duration:", "耗时："))+ui.ValueStyle.Render(r.Duration.Truncate(0).String()),
			ui.LabelStyle.Render(m.text("Status:", "状态："))+ui.SuccessStyle.Render(m.text("success", "成功")),
		)
	}

	for _, line := range infoLines {
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render(m.text("Response:", "返回内容：")))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())

	b.WriteString(ui.RenderHelp([]string{
		m.text("📜 up/down/PgUp/PgDn scroll", "📜 上/下/PgUp/PgDn 滚动"),
		m.text("⏫ Home/End jump", "⏫ Home/End 跳转"),
		m.text("🔁 r retry", "🔁 r 重试"),
		m.text("↩️ b/Esc back", "↩️ b/Esc 返回"),
		m.text("🚪 q quit", "🚪 q 退出"),
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewSettings() string {
	var b strings.Builder

	b.WriteString(m.brandHeader())
	b.WriteString(ui.RenderTitle(m.text("Settings", "设置")))
	if m.config.BaseURL == "" {
		b.WriteString("\n")
		b.WriteString(ui.WarningStyle.Render(m.text(
			"Configure a base URL before using Endpoint TUI. Example: http://127.0.0.1:8080",
			"使用前请配置接口地址。示例：http://127.0.0.1:8080",
		)))
	}
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render(m.text("Base URL:", "接口地址：")))
	b.WriteString("\n")
	b.WriteString(m.settingsBaseURL.View())
	b.WriteString("\n\n")

	b.WriteString(ui.LabelStyle.Render(m.text("Default format:", "默认格式：")))
	b.WriteString("\n")
	for i, enc := range m.encodings {
		if m.settingsOptionCursor == 0 && i == m.settingsEncodingCursor {
			b.WriteString(ui.SelectedStyle.Render("> " + enc))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + enc))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render(m.text("Language:", "语言：")))
	b.WriteString("\n")
	for i, lang := range m.languages {
		label := languageLabel(lang)
		if m.settingsOptionCursor == 1 && i == m.settingsLanguageCursor {
			b.WriteString(ui.SelectedStyle.Render("> " + label))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + label))
		}
		b.WriteString("\n")
	}

	if m.settingsValidationError != "" {
		b.WriteString("\n")
		b.WriteString(ui.ErrorStyle.Render(m.settingsValidationError))
	}

	if m.settingsSaved {
		b.WriteString("\n")
		b.WriteString(ui.SuccessStyle.Render(m.text("Config saved", "配置已保存")))
	}

	b.WriteString(ui.RenderHelp([]string{
		m.text("🔀 Tab switch option", "🔀 Tab 切换选项"),
		m.text("⬆️⬇️ up/down select", "⬆️⬇️ 上/下 选择"),
		m.text("💾 Enter/Ctrl+S save", "💾 Enter/Ctrl+S 保存"),
		m.text("↩️ Esc cancel", "↩️ Esc 取消"),
	}))

	return ui.ContainerStyle.Render(b.String())
}

func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(m.brandHeader())
	b.WriteString(ui.RenderTitle(m.text("Endpoint TUI", "60 秒新闻终端")))
	b.WriteString("\n")

	if m.loadErr != nil {
		b.WriteString(ui.ErrorStyle.Render(m.text("Load failed: ", "加载失败：")))
		b.WriteString(ui.ValueStyle.Render(safeErrorMessage(m.loadErr)))
	}

	b.WriteString(ui.RenderHelp([]string{
		m.text("🔄 r reload", "🔄 r 重新加载"),
		m.text("⚙️ s settings", "⚙️ s 设置"),
		m.text("🚪 q quit", "🚪 q 退出"),
	}))

	return ui.ContainerStyle.Render(b.String())
}

func formatEndpointLine(number int, ep api.Endpoint, width int) string {
	path := ep.Path
	name := ep.Name
	prefix := fmt.Sprintf("%d. ", number)

	if name == path {
		return ui.Truncate(prefix+path, width)
	}

	line := prefix + ui.PadRight(name, 20) + " " + path
	return ui.Truncate(line, width)
}

func (m Model) text(en, zh string) string {
	if m.config.Language == "zh" {
		return zh
	}
	return en
}

func (m Model) brandHeader() string {
	var b strings.Builder
	b.WriteString(ui.WarningStyle.Render(m.text(
		"📰 Understand the world in 60 seconds",
		"📰 读懂世界 · 每天 60 秒读懂世界",
	)))
	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render(m.text(
		"✨ Daily curated news for major world events",
		"✨ 获取每日精选新闻，快速了解世界大事",
	)))
	b.WriteString("\n")
	b.WriteString(ui.RenderTitle(m.text("🌍 Endpoint TUI", "🌍 60 秒新闻终端")))
	b.WriteString("\n\n")
	return b.String()
}

func languageLabel(lang string) string {
	switch lang {
	case "zh":
		return "中文"
	default:
		return "English"
	}
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
