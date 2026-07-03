// Package config manages the application config file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config is the persisted application configuration.
type Config struct {
	BaseURL         string `json:"base_url"`
	DefaultEncoding string `json:"default_encoding"`
	Language        string `json:"language"`
}

// DefaultConfig returns default settings.
func DefaultConfig() Config {
	return Config{
		BaseURL:         "",
		DefaultEncoding: "json",
		Language:        "en",
	}
}

var validEncodings = map[string]bool{
	"json":     true,
	"text":     true,
	"markdown": true,
}

var validLanguages = map[string]bool{
	"en": true,
	"zh": true,
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".config", "endpoint-tui"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func ensureDir() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}
	return nil
}

// Load reads the default config file.
func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to determine config path: %w", err)
	}
	return LoadFromPath(path)
}

// LoadFromPath reads config from a specific path.
func LoadFromPath(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("invalid config JSON: %w", err)
	}

	if !validEncodings[cfg.DefaultEncoding] {
		cfg.DefaultEncoding = "json"
	}
	if !validLanguages[cfg.Language] {
		cfg.Language = "en"
	}

	return cfg, nil
}

// Save writes the config file.
func Save(cfg Config) error {
	if !validEncodings[cfg.DefaultEncoding] {
		cfg.DefaultEncoding = "json"
	}
	if !validLanguages[cfg.Language] {
		cfg.Language = "en"
	}

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
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// ValidateBaseURL validates and normalizes the base URL.
func ValidateBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("base URL cannot be empty")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "", errors.New("base URL must start with http:// or https://")
	}
	raw = strings.TrimRight(raw, "/")
	return raw, nil
}
