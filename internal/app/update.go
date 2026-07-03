package app

import (
	"context"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/config"
	"endpoint-tui/internal/urlutil"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchEndpointsCmd(url string) tea.Cmd {
	return func() tea.Msg {
		endpoints, err := api.FetchEndpoints(url)
		return EndpointsLoadedMsg{
			Endpoints: endpoints,
			Error:     err,
		}
	}
}

func executeCurlCmd(requestURL string, executor api.CommandExecutor) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		result := executor.Execute(ctx, requestURL)
		return CurlResultMsg{Result: result}
	}
}

// Update handles Bubble Tea messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 12
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.page == PageEndpointList || m.page == PageError {
				return m, tea.Quit
			}
		}

	case EndpointsLoadedMsg:
		m.loading = false
		if msg.Error != nil {
			m.loadErr = msg.Error
			m.page = PageError
		} else {
			m.endpoints = msg.Endpoints
			m.cursor = 0
			m.loadErr = nil
			if m.config.BaseURL == "" {
				m.page = PageSettings
			} else {
				m.page = PageEndpointList
			}
		}
		return m, nil

	case CurlResultMsg:
		m.result = msg.Result
		m.viewport.SetContent(formatResultContent(msg.Result, m.SelectedEncoding()))
		m.viewport.GotoTop()
		m.page = PageResult
		return m, nil
	}

	switch m.page {
	case PageLoading:
		return m.updateLoading(msg)
	case PageEndpointList:
		return m.updateEndpointList(msg)
	case PageEncodingSelect:
		return m.updateEncodingSelect(msg)
	case PageResult:
		return m.updateResult(msg)
	case PageSettings:
		return m.updateSettings(msg)
	case PageError:
		return m.updateError(msg)
	default:
		return m, nil
	}
}

func (m Model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) updateEndpointList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.endpoints)-1 {
				m.cursor++
			}
		case "enter":
			if m.SelectedEndpoint() != nil {
				m.page = PageEncodingSelect
			}
		case "r":
			m.loading = true
			m.page = PageLoading
			return m, fetchEndpointsCmd(m.discoveryURL)
		case "s":
			m.settingsBaseURL.SetValue(m.config.BaseURL)
			encIdx := 0
			for i, e := range m.encodings {
				if e == m.config.DefaultEncoding {
					encIdx = i
					break
				}
			}
			m.settingsEncodingCursor = encIdx
			m.settingsSaved = false
			m.settingsValidationError = ""
			m.page = PageSettings
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateEncodingSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.encodingCursor > 0 {
				m.encodingCursor--
			}
		case "down", "j":
			if m.encodingCursor < len(m.encodings)-1 {
				m.encodingCursor++
			}
		case "enter":
			ep := m.SelectedEndpoint()
			if ep == nil {
				return m, nil
			}
			requestURL, err := urlutil.BuildURL(m.config.BaseURL, ep.Path, m.SelectedEncoding())
			if err != nil {
				m.loadErr = err
				m.page = PageError
				return m, nil
			}
			m.result = api.CurlResult{}
			executor := api.NewCurlExecutor()
			return m, executeCurlCmd(requestURL, executor)
		case "esc":
			m.page = PageEndpointList
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup":
			m.viewport.PageUp()
		case "pgdown":
			m.viewport.PageDown()
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		case "r":
			ep := m.SelectedEndpoint()
			if ep == nil {
				return m, nil
			}
			requestURL, err := urlutil.BuildURL(m.config.BaseURL, ep.Path, m.SelectedEncoding())
			if err != nil {
				m.loadErr = err
				m.page = PageError
				return m, nil
			}
			executor := api.NewCurlExecutor()
			return m, executeCurlCmd(requestURL, executor)
		case "b", "esc":
			m.page = PageEndpointList
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.page = PageEndpointList
			return m, nil
		case "ctrl+s":
			baseURL, err := config.ValidateBaseURL(m.settingsBaseURL.Value())
			if err != nil {
				m.settingsValidationError = err.Error()
				return m, nil
			}
			m.config.BaseURL = baseURL
			m.config.DefaultEncoding = m.SettingsSelectedEncoding()
			if err := config.Save(m.config); err != nil {
				m.settingsValidationError = "save failed: " + err.Error()
				return m, nil
			}
			m.settingsSaved = true
			m.settingsValidationError = ""
			for i, e := range m.encodings {
				if e == m.config.DefaultEncoding {
					m.encodingCursor = i
					break
				}
			}
			return m, nil
		case "up", "k":
			if m.settingsEncodingCursor > 0 {
				m.settingsEncodingCursor--
			}
		case "down", "j":
			if m.settingsEncodingCursor < len(m.encodings)-1 {
				m.settingsEncodingCursor++
			}
		default:
			m.settingsSaved = false
			m.settingsValidationError = ""
		}
	}

	m.settingsBaseURL, cmd = m.settingsBaseURL.Update(msg)
	return m, cmd
}

func (m Model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			m.loadErr = nil
			m.page = PageLoading
			return m, fetchEndpointsCmd(m.discoveryURL)
		case "s":
			m.settingsBaseURL.SetValue(m.config.BaseURL)
			m.settingsSaved = false
			m.settingsValidationError = ""
			m.page = PageSettings
			return m, nil
		}
	}
	return m, nil
}
