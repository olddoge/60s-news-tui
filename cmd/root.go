// Package cmd 处理命令行参数解析。
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"endpoint-tui/internal/api"
	"endpoint-tui/internal/config"
)

// Version 是程序版本号，编译时可通过 ldflags 注入。
var Version = "0.1.0"

// Parse 解析命令行参数并返回配置和发现服务地址。
func Parse() (config.Config, string) {
	var configPath string
	var discoveryURL string
	var showVersion bool
	var showHelp bool

	defaultConfigPath := ""
	if home, err := os.UserHomeDir(); err == nil {
		defaultConfigPath = filepath.Join(home, ".config", "endpoint-tui", "config.json")
	}

	flag.StringVar(&configPath, "config", defaultConfigPath, "配置文件路径")
	flag.StringVar(&discoveryURL, "discovery-url", api.GetDiscoveryURL(), "接口发现服务地址")
	flag.BoolVar(&showVersion, "version", false, "显示版本号")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Endpoint TUI - 终端接口调试工具

用法:
  endpoint-tui [选项]

选项:
  --config FILE        配置文件路径 (默认: ~/.config/endpoint-tui/config.json)
  --discovery-url URL  接口发现服务地址 (默认: %s)
  --version            显示版本号
  --help               显示帮助信息
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

	// 加载配置文件
	cfg, err := config.LoadFromPath(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
		cfg = config.DefaultConfig()
	}

	return cfg, discoveryURL
}
