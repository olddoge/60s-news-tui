// Package ui 管理 TUI 的各个页面组件和样式定义。
package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	// 基础颜色
	primaryColor   = lipgloss.Color("39")  // 蓝色
	successColor   = lipgloss.Color("42")  // 绿色
	errorColor     = lipgloss.Color("196") // 红色
	warningColor   = lipgloss.Color("220") // 黄色
	mutedColor     = lipgloss.Color("245") // 灰色
	dimColor       = lipgloss.Color("240") // 暗灰色
	bgColor        = lipgloss.Color("235") // 暗背景
	highlightColor = lipgloss.Color("63")  // 高亮蓝

	noColor bool
)

func init() {
	// 检查 NO_COLOR 环境变量
	if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}
}

// color 在 NO_COLOR 模式下返回空样式，否则返回带颜色样式。
func color(c lipgloss.Color) lipgloss.TerminalColor {
	if noColor {
		return lipgloss.NoColor{}
	}
	return c
}

// 全局样式定义

// TitleStyle 程序标题样式。
var TitleStyle = lipgloss.NewStyle().
	Foreground(color(primaryColor)).
	Bold(true).
	MarginBottom(1)

// SubtitleStyle 副标题样式。
var SubtitleStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor)).
	MarginBottom(1)

// SelectedStyle 选中项样式。
var SelectedStyle = lipgloss.NewStyle().
	Foreground(color(highlightColor)).
	Bold(true).
	PaddingLeft(2)

// NormalStyle 未选中项样式。
var NormalStyle = lipgloss.NewStyle().
	Foreground(color(lipgloss.Color("252"))).
	PaddingLeft(2)

// ErrorStyle 错误信息样式。
var ErrorStyle = lipgloss.NewStyle().
	Foreground(color(errorColor)).
	Bold(true)

// SuccessStyle 成功信息样式。
var SuccessStyle = lipgloss.NewStyle().
	Foreground(color(successColor))

// HelpStyle 底部快捷键说明样式。
var HelpStyle = lipgloss.NewStyle().
	Foreground(color(dimColor)).
	MarginTop(1)

// StatusStyle 状态信息样式。
var StatusStyle = lipgloss.NewStyle().
	Foreground(color(warningColor))

// WarningStyle 警告信息样式。
var WarningStyle = lipgloss.NewStyle().
	Foreground(color(warningColor)).
	Bold(true)

// LoadingStyle 加载中样式。
var LoadingStyle = lipgloss.NewStyle().
	Foreground(color(primaryColor)).
	Italic(true)

// InfoStyle 信息标签样式。
var InfoStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor))

// LabelStyle 标签样式。
var LabelStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor)).
	Width(10)

// ValueStyle 值样式。
var ValueStyle = lipgloss.NewStyle().
	Foreground(color(lipgloss.Color("252")))

// InputStyle 输入框样式。
var InputStyle = lipgloss.NewStyle().
	Foreground(color(highlightColor)).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	BorderForeground(color(primaryColor)).
	Width(60)

// ContainerStyle 容器样式。
var ContainerStyle = lipgloss.NewStyle().
	Padding(1, 2)

// 布局工具

// AppWidth 是应用程序的标准内容宽度。
const AppWidth = 80

// RenderTitle 渲染带标题的页面头部。
func RenderTitle(title string) string {
	return TitleStyle.Render(title)
}

// RenderHelp 渲染底部帮助栏。
func RenderHelp(keys []string) string {
	help := ""
	for i, k := range keys {
		if i > 0 {
			help += "  "
		}
		help += HelpStyle.Render(k)
	}
	return "\n" + help
}

// Truncate 截断字符串到最大长度。
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// PadRight 右侧填充到指定宽度。
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + spaces(width-len(s))
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
