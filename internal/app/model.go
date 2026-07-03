package app

import (
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

	encodings      []string
	encodingCursor int

	result api.CurlResult

	config       config.Config
	discoveryURL string

	settingsBaseURL         textinput.Model
	settingsEncodingCursor  int
	settingsSaved           bool
	settingsValidationError string

	viewport viewport.Model

	width  int
	height int
	ready  bool

	loading bool
	loadErr error
}

// Encodings is the selectable response encoding list.
var Encodings = []string{"json", "text", "markdown"}

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

	return Model{
		page:                   PageLoading,
		encodings:              Encodings,
		encodingCursor:         encIdx,
		config:                 cfg,
		discoveryURL:           discoveryURL,
		settingsBaseURL:        ti,
		settingsEncodingCursor: encIdx,
		loading:                true,
	}
}

// SelectedEndpoint returns the currently selected endpoint.
func (m Model) SelectedEndpoint() *api.Endpoint {
	if m.cursor < 0 || m.cursor >= len(m.endpoints) {
		return nil
	}
	return &m.endpoints[m.cursor]
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

// Init starts loading endpoints.
func (m Model) Init() tea.Cmd {
	return fetchEndpointsCmd(m.discoveryURL)
}
