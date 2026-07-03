package app

import (
	"strings"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the top-level Bubble Tea model.
type Model struct {
	page Page

	endpoints []api.Endpoint
	cursor    int
	searching bool
	search    string

	encodings      []string
	encodingCursor int
	languages      []string
	languageCursor int

	result api.CurlResult

	config       config.Config
	discoveryURL string

	settingsBaseURL         textinput.Model
	settingsEncodingCursor  int
	settingsLanguageCursor  int
	settingsOptionCursor    int
	settingsSaved           bool
	settingsValidationError string

	viewport viewport.Model

	width  int
	height int
	ready  bool

	loading        bool
	loadingMessage string
	loadErr        error
}

// Encodings is the selectable response encoding list.
var Encodings = []string{"json", "text", "markdown"}

// Languages is the selectable UI language list.
var Languages = []string{"en", "zh"}

// NewModel creates the application model.
func NewModel(cfg config.Config, discoveryURL string) Model {
	ti := textinput.New()
	ti.Placeholder = "for example: http://127.0.0.1:8080"
	ti.CharLimit = 256
	ti.Width = 60
	ti.SetValue(cfg.BaseURL)
	ti.Focus()

	encIdx := 0
	for i, e := range Encodings {
		if e == cfg.DefaultEncoding {
			encIdx = i
			break
		}
	}
	langIdx := 0
	for i, lang := range Languages {
		if lang == cfg.Language {
			langIdx = i
			break
		}
	}

	page := PageLoading
	loading := true
	loadingMessage := "Loading endpoint list..."
	if cfg.BaseURL == "" {
		page = PageSettings
		loading = false
		loadingMessage = ""
	}

	return Model{
		page:                   page,
		encodings:              Encodings,
		encodingCursor:         encIdx,
		languages:              Languages,
		languageCursor:         langIdx,
		config:                 cfg,
		discoveryURL:           discoveryURL,
		settingsBaseURL:        ti,
		settingsEncodingCursor: encIdx,
		settingsLanguageCursor: langIdx,
		loading:                loading,
		loadingMessage:         loadingMessage,
	}
}

// SelectedEndpoint returns the currently selected endpoint.
func (m Model) SelectedEndpoint() *api.Endpoint {
	indexes := m.filteredEndpointIndexes()
	if m.cursor < 0 || m.cursor >= len(indexes) {
		return nil
	}
	return &m.endpoints[indexes[m.cursor]]
}

// SelectedEncoding returns the currently selected response encoding.
func (m Model) SelectedEncoding() string {
	if m.encodingCursor < 0 || m.encodingCursor >= len(m.encodings) {
		return "json"
	}
	return m.encodings[m.encodingCursor]
}

// SettingsSelectedEncoding returns the encoding selected in settings.
func (m Model) SettingsSelectedEncoding() string {
	if m.settingsEncodingCursor < 0 || m.settingsEncodingCursor >= len(m.encodings) {
		return "json"
	}
	return m.encodings[m.settingsEncodingCursor]
}

func (m Model) SelectedLanguage() string {
	if m.languageCursor < 0 || m.languageCursor >= len(m.languages) {
		return "en"
	}
	return m.languages[m.languageCursor]
}

func (m Model) SettingsSelectedLanguage() string {
	if m.settingsLanguageCursor < 0 || m.settingsLanguageCursor >= len(m.languages) {
		return "en"
	}
	return m.languages[m.settingsLanguageCursor]
}

func (m Model) filteredEndpointIndexes() []int {
	query := strings.TrimSpace(strings.ToLower(m.search))
	indexes := make([]int, 0, len(m.endpoints))
	for i, ep := range m.endpoints {
		if query == "" ||
			strings.Contains(strings.ToLower(ep.Name), query) ||
			strings.Contains(strings.ToLower(ep.Path), query) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (m Model) endpointDiscoveryURL() string {
	if m.config.BaseURL != "" && (m.discoveryURL == "" || m.discoveryURL == api.DefaultDiscoveryURL) {
		return m.config.BaseURL
	}
	return m.discoveryURL
}

// EndpointDiscoveryURL returns the discovery URL used to load endpoint lists.
func (m Model) EndpointDiscoveryURL() string {
	return m.endpointDiscoveryURL()
}

// Init starts loading endpoints.
func (m Model) Init() tea.Cmd {
	if m.config.BaseURL == "" {
		return nil
	}
	return fetchEndpointsCmd(m.endpointDiscoveryURL())
}
