package app

import (
	"endpoint-tui/internal/api"
	"endpoint-tui/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model 是 Bubble Tea 的顶层 Model，管理所有页面状态。
type Model struct {
	// 当前页面
	page Page

	// 接口数据
	endpoints []api.Endpoint
	cursor    int // 接口列表选中索引

	// encoding 选择
	encodings      []string
	encodingCursor int

	// 请求结果
	result api.CurlResult

	// 配置
	config       config.Config
	discoveryURL string

	// 设置页
	settingsBaseURL         textinput.Model
	settingsEncodingCursor  int
	settingsSaved           bool
	settingsValidationError string

	// 结果页
	viewport viewport.Model

	// 终端尺寸
	width  int
	height int
	ready  bool

	// 加载状态
	loading bool
	loadErr error
}

// Encodings 是可选的 encoding 值列表。
var Encodings = []string{"json", "text", "markdown"}

// NewModel 创建应用 Model。
func NewModel(cfg config.Config, discoveryURL string) Model {
	// 设置页输入框
	ti := textinput.New()
	ti.Placeholder = "例如：http://127.0.0.1:8080"
	ti.CharLimit = 256
	ti.Width = 60
	ti.SetValue(cfg.BaseURL)
	ti.Focus()

	// 确定默认 encoding 位置
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

// SelectedEndpoint 返回当前选中的接口。
func (m Model) SelectedEndpoint() *api.Endpoint {
	if m.cursor < 0 || m.cursor >= len(m.endpoints) {
		return nil
	}
	return &m.endpoints[m.cursor]
}

// SelectedEncoding 返回当前选中的 encoding。
func (m Model) SelectedEncoding() string {
	if m.encodingCursor < 0 || m.encodingCursor >= len(m.encodings) {
		return "json"
	}
	return m.encodings[m.encodingCursor]
}

// SettingsSelectedEncoding 返回设置页当前选中的 encoding。
func (m Model) SettingsSelectedEncoding() string {
	if m.settingsEncodingCursor < 0 || m.settingsEncodingCursor >= len(m.encodings) {
		return "json"
	}
	return m.encodings[m.settingsEncodingCursor]
}

// Init 是 Bubble Tea 的初始化命令，启动时自动获取接口列表。
func (m Model) Init() tea.Cmd {
	return fetchEndpointsCmd(m.discoveryURL)
}
