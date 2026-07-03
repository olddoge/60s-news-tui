// Package ui defines shared TUI styles and layout helpers.
package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("39")
	successColor   = lipgloss.Color("42")
	errorColor     = lipgloss.Color("196")
	warningColor   = lipgloss.Color("220")
	mutedColor     = lipgloss.Color("245")
	dimColor       = lipgloss.Color("240")
	highlightColor = lipgloss.Color("63")
	noColor        bool
)

func init() {
	if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}
}

func color(c lipgloss.Color) lipgloss.TerminalColor {
	if noColor {
		return lipgloss.NoColor{}
	}
	return c
}

var TitleStyle = lipgloss.NewStyle().
	Foreground(color(primaryColor)).
	Bold(true).
	MarginBottom(1)

var SubtitleStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor)).
	MarginBottom(1)

var SelectedStyle = lipgloss.NewStyle().
	Foreground(color(highlightColor)).
	Bold(true).
	PaddingLeft(2)

var NormalStyle = lipgloss.NewStyle().
	Foreground(color(lipgloss.Color("252"))).
	PaddingLeft(2)

var ErrorStyle = lipgloss.NewStyle().
	Foreground(color(errorColor)).
	Bold(true)

var SuccessStyle = lipgloss.NewStyle().
	Foreground(color(successColor))

var HelpStyle = lipgloss.NewStyle().
	Foreground(color(dimColor)).
	MarginTop(1)

var StatusStyle = lipgloss.NewStyle().
	Foreground(color(warningColor))

var WarningStyle = lipgloss.NewStyle().
	Foreground(color(warningColor)).
	Bold(true)

var LoadingStyle = lipgloss.NewStyle().
	Foreground(color(primaryColor)).
	Italic(true)

var InfoStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor))

var LabelStyle = lipgloss.NewStyle().
	Foreground(color(mutedColor)).
	Width(10)

var ValueStyle = lipgloss.NewStyle().
	Foreground(color(lipgloss.Color("252")))

var InputStyle = lipgloss.NewStyle().
	Foreground(color(highlightColor)).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	BorderForeground(color(primaryColor)).
	Width(60)

var ContainerStyle = lipgloss.NewStyle().
	Padding(1, 2)

const AppWidth = 80

func RenderTitle(title string) string {
	return TitleStyle.Render(title)
}

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

func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

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
