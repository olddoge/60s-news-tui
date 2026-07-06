package main

import (
	"fmt"
	"os"

	"endpoint-tui/cmd"
	"endpoint-tui/internal/api"
	"endpoint-tui/internal/app"
	"endpoint-tui/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, discoveryURL := cmd.Parse()
	config.SetDefaultPublicInstances(publicInstanceJSON)

	if err := api.CheckCurlAvailable(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	m := app.NewModel(cfg, discoveryURL)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "program failed: %v\n", err)
		os.Exit(1)
	}
}
