// Package cmd handles command-line arguments.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/config"
)

// Version is injected at build time with ldflags.
var Version = "0.1.0"

// Parse parses command-line arguments and returns config plus discovery URL.
func Parse() (config.Config, string) {
	var configPath string
	var discoveryURL string
	var showVersion bool
	var showHelp bool

	defaultConfigPath := ""
	if home, err := os.UserHomeDir(); err == nil {
		defaultConfigPath = filepath.Join(home, ".config", "endpoint-tui", "config.json")
	}

	flag.StringVar(&configPath, "config", defaultConfigPath, "config file path")
	flag.StringVar(&discoveryURL, "discovery-url", api.GetDiscoveryURL(), "endpoint discovery service URL")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showHelp, "help", false, "show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Endpoint TUI - terminal API browser

Usage:
  endpoint-tui [options]

Options:
  --config FILE        config file path (default: ~/.config/endpoint-tui/config.json)
  --discovery-url URL  endpoint discovery service URL (default: %s)
  --version            show version
  --help               show help
`, api.GetDiscoveryURL())
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("endpoint-tui version %s\n", Version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	cfg, err := config.LoadFromPath(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		cfg = config.DefaultConfig()
	}

	return cfg, discoveryURL
}
