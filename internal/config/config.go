// Package config 管理程序的配置文件读写。
// 配置文件位于 ~/.config/endpoint-tui/config.json，权限为 0600。
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 表示程序配置结构。
type Config struct {
	BaseURL         string `json:"base_url"`
	DefaultEncoding string `json:"default_encoding"`
}

// DefaultConfig 返回默认配置。
func DefaultConfig() Config {
	return Config{
		BaseURL:         "",
		DefaultEncoding: "json",
	}
}

// validEncodings 是允许的 encoding 值集合。
var validEncodings = map[string]bool{
	"json":     true,
	"text":     true,
	"markdown": true,
}

// configDir 返回配置文件所在目录。
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}
	return filepath.Join(home, ".config", "endpoint-tui"), nil
}

// configPath 返回配置文件的完整路径。
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// ensureDir 确保配置目录存在，不存在则创建。
func ensureDir() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("无法创建配置目录 %s: %w", dir, err)
	}
	return nil
}

// Load 从默认配置文件加载配置。文件不存在时返回默认配置。
func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), fmt.Errorf("无法确定配置文件路径，使用默认配置: %w", err)
	}
	return LoadFromPath(path)
}

// LoadFromPath 从指定路径加载配置。文件不存在时返回默认配置。
func LoadFromPath(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("无法读取配置文件 %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("配置文件 JSON 格式错误: %w", err)
	}

	// 校验 encoding 值
	if !validEncodings[cfg.DefaultEncoding] {
		cfg.DefaultEncoding = "json"
	}

	return cfg, nil
}

// Save 保存配置到文件。
func Save(cfg Config) error {
	// 校验 encoding
	if !validEncodings[cfg.DefaultEncoding] {
		cfg.DefaultEncoding = "json"
	}

	// 清理根路径
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	if err := ensureDir(); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("无法序列化配置: %w", err)
	}

	// 写入临时文件后原子重命名
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("无法写入配置文件: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("无法保存配置文件: %w", err)
	}

	return nil
}

// ValidateBaseURL 校验根路径格式。返回清理后的 URL 和可能的错误。
func ValidateBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("根路径不能为空")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "", errors.New("根路径必须以 http:// 或 https:// 开头")
	}
	raw = strings.TrimRight(raw, "/")
	return raw, nil
}
