package app

import (
	"context"
	"unicode/utf8"

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

func (m Model) startSelectedRequest() (Model, tea.Cmd) {
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
	m.loading = true
	m.loadingMessage = "Requesting endpoint..."
	m.page = PageLoading
	executor := api.NewCurlExecutor()
	return m, executeCurlCmd(requestURL, executor)
}

func (m Model) saveSettings() (Model, tea.Cmd) {
	baseURL, err := config.ValidateBaseURL(m.settingsBaseURL.Value())
	if err != nil {
		m.settingsValidationError = err.Error()
		return m, nil
	}
	m.config.BaseURL = baseURL
	m.config.DefaultEncoding = m.SettingsSelectedEncoding()
	m.config.Language = m.SettingsSelectedLanguage()
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
	for i, lang := range m.languages {
		if lang == m.config.Language {
			m.languageCursor = i
			break
		}
	}
	m.loading = true
	m.loadingMessage = "Config saved. Loading endpoint list..."
	m.page = PageLoading
	return m, fetchEndpointsCmd(m.endpointDiscoveryURL())
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
			if m.page == PageEndpointList && m.searching {
				break
			}
			return m, tea.Quit
		}

	case EndpointsLoadedMsg:
		m.loading = false
		m.loadingMessage = ""
		if msg.Error != nil {
			m.loadErr = msg.Error
			m.page = PageError
		} else {
			m.endpoints = api.LocalizeEndpointsForLanguage(msg.Endpoints, m.config.Language)
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
		m.loading = false
		m.loadingMessage = ""
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
		if m.searching {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.filteredEndpointIndexes())-1 {
					m.cursor++
				}
			case "esc":
				m.searching = false
				m.search = ""
				m.cursor = 0
			case "enter":
				m.searching = false
			case "backspace":
				if m.search != "" {
					_, size := utf8.DecodeLastRuneInString(m.search)
					m.search = m.search[:len(m.search)-size]
					m.cursor = 0
				}
			default:
				if len(msg.Runes) > 0 {
					m.search += string(msg.Runes)
					m.cursor = 0
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filteredEndpointIndexes())-1 {
				m.cursor++
			}
		case "enter":
			return m.startSelectedRequest()
		case "/":
			m.searching = true
		case "r":
			m.loading = true
			m.loadingMessage = "Loading endpoint list..."
			m.page = PageLoading
			return m, fetchEndpointsCmd(m.endpointDiscoveryURL())
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
			langIdx := 0
			for i, lang := range m.languages {
				if lang == m.config.Language {
					langIdx = i
					break
				}
			}
			m.settingsLanguageCursor = langIdx
			m.settingsOptionCursor = 0
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
			return m.startSelectedRequest()
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
			return m.startSelectedRequest()
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
			if m.config.BaseURL == "" {
				m.settingsValidationError = "Base URL is required before using Endpoint TUI"
				return m, nil
			}
			m.page = PageEndpointList
			return m, nil
		case "ctrl+s", "enter":
			return m.saveSettings()
		case "tab":
			m.settingsOptionCursor = (m.settingsOptionCursor + 1) % 2
		case "up", "k":
			if m.settingsOptionCursor == 0 {
				if m.settingsEncodingCursor > 0 {
					m.settingsEncodingCursor--
				}
			} else if m.settingsLanguageCursor > 0 {
				m.settingsLanguageCursor--
			}
		case "down", "j":
			if m.settingsOptionCursor == 0 {
				if m.settingsEncodingCursor < len(m.encodings)-1 {
					m.settingsEncodingCursor++
				}
			} else if m.settingsLanguageCursor < len(m.languages)-1 {
				m.settingsLanguageCursor++
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
			m.loadingMessage = "Loading endpoint list..."
			m.loadErr = nil
			m.page = PageLoading
			return m, fetchEndpointsCmd(m.endpointDiscoveryURL())
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
