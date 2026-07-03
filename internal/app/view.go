package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/ui"
)

// View 是 Bubble Tea 的主渲染函数。
func (m Model) View() string {
	if !m.ready {
		return "正在初始化..."
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

// viewLoading 渲染加载页面。
func (m Model) viewLoading() string {
	return ui.ContainerStyle.Render(
		ui.LoadingStyle.Render("正在加载接口列表..."),
	)
}

// viewEndpointList 渲染接口列表页面。
func (m Model) viewEndpointList() string {
	var b strings.Builder

	// 标题
	b.WriteString(ui.RenderTitle("Endpoint TUI"))
	b.WriteString("\n")

	// 根路径
	baseURL := m.config.BaseURL
	if baseURL == "" {
		baseURL = "(未设置)"
	}
	b.WriteString(ui.InfoStyle.Render("根路径：" + baseURL))
	b.WriteString("\n\n")

	// 接口列表
	b.WriteString(ui.InfoStyle.Render("接口列表："))
	b.WriteString("\n\n")

	// 计算可视区域
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

	// 滚动指示器
	if len(m.endpoints) > listHeight {
		b.WriteString(ui.HelpStyle.Render(
			fmt.Sprintf("  ... %d/%d", m.cursor+1, len(m.endpoints)),
		))
		b.WriteString("\n")
	}

	// 帮助栏
	b.WriteString(ui.RenderHelp([]string{
		"↑/↓ 选择",
		"Enter 请求",
		"s 设置",
		"r 刷新",
		"q 退出",
	}))

	return ui.ContainerStyle.Render(b.String())
}

// viewEncodingSelect 渲染 encoding 选择页面。
func (m Model) viewEncodingSelect() string {
	var b strings.Builder

	ep := m.SelectedEndpoint()
	epPath := ""
	if ep != nil {
		epPath = ep.Path
	}

	b.WriteString(ui.RenderTitle("请求配置"))
	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render("接口：" + epPath))
	b.WriteString("\n\n")
	b.WriteString(ui.InfoStyle.Render("请选择返回格式："))
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
		"↑/↓ 选择",
		"Enter 执行请求",
		"Esc 返回列表",
	}))

	return ui.ContainerStyle.Render(b.String())
}

// viewResult 渲染请求结果页面。
func (m Model) viewResult() string {
	var b strings.Builder

	r := m.result
	ep := m.SelectedEndpoint()
	epPath := ""
	if ep != nil {
		epPath = ep.Name
	}

	b.WriteString(ui.RenderTitle("请求完成"))
	b.WriteString("\n")

	// 信息行
	infoLines := []string{
		ui.LabelStyle.Render("接口：") + ui.ValueStyle.Render(epPath),
		ui.LabelStyle.Render("格式：") + ui.ValueStyle.Render(m.SelectedEncoding()),
		ui.LabelStyle.Render("地址：") + ui.ValueStyle.Render(ui.Truncate(r.URL, m.width-10)),
	}

	if r.Cancelled {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("状态：")+ui.WarningStyle.Render("已取消"),
		)
	} else if r.Error != nil || r.ExitCode != 0 {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("状态：")+ui.ErrorStyle.Render(fmt.Sprintf("失败 (exit=%d)", r.ExitCode)),
		)
		if r.Stderr != "" {
			infoLines = append(infoLines,
				ui.LabelStyle.Render("错误：")+ui.ErrorStyle.Render(ui.Truncate(r.Stderr, m.width-10)),
			)
		}
	} else {
		infoLines = append(infoLines,
			ui.LabelStyle.Render("耗时：")+ui.ValueStyle.Render(r.Duration.Truncate(0).String()),
			ui.LabelStyle.Render("状态：")+ui.SuccessStyle.Render("成功"),
		)
	}

	for _, line := range infoLines {
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.InfoStyle.Render("返回内容："))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())

	// 帮助栏
	b.WriteString(ui.RenderHelp([]string{
		"↑/↓/PgUp/PgDn 滚动",
		"Home/End 跳转",
		"r 重新请求",
		"b/Esc 返回列表",
		"q 退出",
	}))

	return ui.ContainerStyle.Render(b.String())
}

// viewSettings 渲染设置页面。
func (m Model) viewSettings() string {
	var b strings.Builder

	b.WriteString(ui.RenderTitle("设置"))
	b.WriteString("\n")

	// 根路径
	b.WriteString(ui.LabelStyle.Render("根路径："))
	b.WriteString("\n")
	b.WriteString(m.settingsBaseURL.View())
	b.WriteString("\n\n")

	// 默认 encoding
	b.WriteString(ui.LabelStyle.Render("默认格式："))
	b.WriteString("\n")
	for i, enc := range m.encodings {
		if i == m.settingsEncodingCursor {
			b.WriteString(ui.SelectedStyle.Render("> " + enc))
		} else {
			b.WriteString(ui.NormalStyle.Render("  " + enc))
		}
		b.WriteString("\n")
	}

	// 验证错误
	if m.settingsValidationError != "" {
		b.WriteString("\n")
		b.WriteString(ui.ErrorStyle.Render(m.settingsValidationError))
	}

	// 保存状态
	if m.settingsSaved {
		b.WriteString("\n")
		b.WriteString(ui.SuccessStyle.Render("配置已保存"))
	}

	b.WriteString(ui.RenderHelp([]string{
		"↑/↓ 选择格式",
		"Ctrl+S 保存",
		"Esc 取消",
	}))

	return ui.ContainerStyle.Render(b.String())
}

// viewError 渲染错误页面。
func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(ui.RenderTitle("Endpoint TUI"))
	b.WriteString("\n")

	if m.loadErr != nil {
		b.WriteString(ui.ErrorStyle.Render("加载失败："))
		b.WriteString(ui.ValueStyle.Render(m.loadErr.Error()))
	}

	b.WriteString(ui.RenderHelp([]string{
		"r 重新加载",
		"s 设置",
		"q 退出",
	}))

	return ui.ContainerStyle.Render(b.String())
}

// formatEndpointLine 格式化单行接口显示。
func formatEndpointLine(ep api.Endpoint, width int) string {
	path := ep.Path
	name := ep.Name

	// 如果名称和路径相同，只显示路径
	if name == path {
		return ui.Truncate(path, width)
	}

	// 名称 + 路径
	line := ui.PadRight(name, 20) + " " + path
	return ui.Truncate(line, width)
}

// formatResultContent 格式化请求结果内容。
func formatResultContent(r api.CurlResult, encoding string) string {
	if r.Cancelled {
		return "请求已取消。"
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
		return "(返回内容为空)"
	}

	// 根据 encoding 处理内容
	switch encoding {
	case "json":
		return formatJSON(content)
	case "markdown", "text":
		return sanitizeContent(content)
	default:
		return sanitizeContent(content)
	}
}

// formatJSON 格式化 JSON 内容，失败时返回原文。
func formatJSON(content string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(content), &obj); err != nil {
		// 不是有效 JSON，返回原文
		return sanitizeContent(content)
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return sanitizeContent(content)
	}
	return string(formatted)
}

// sanitizeContent 过滤控制字符，防止终端注入。
func sanitizeContent(content string) string {
	var b strings.Builder
	b.Grow(len(content))

	for _, r := range content {
		// 保留换行、制表符和可打印字符
		if r == '\n' || r == '\t' || r == '\r' {
			b.WriteRune(r)
			continue
		}
		if r < 32 || (r >= 0x7F && r <= 0x9F) {
			// 跳过控制字符（ansi 转义序列等）
			continue
		}
		if !utf8.ValidRune(r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
