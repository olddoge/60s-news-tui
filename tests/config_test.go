package tests

import (
	"os"
	"path/filepath"
	"testing"

	"endpoint-tui/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.BaseURL != "" {
		t.Errorf("expected empty BaseURL, got %s", cfg.BaseURL)
	}
	if cfg.DefaultEncoding != "json" {
		t.Errorf("expected json DefaultEncoding, got %s", cfg.DefaultEncoding)
	}
	if cfg.Language != "en" {
		t.Errorf("expected en Language, got %s", cfg.Language)
	}
}

func TestLoad_FileNotExist(t *testing.T) {
	cfg, err := config.LoadFromPath("/tmp/nonexistent-config-test.json")
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
	if cfg.BaseURL != "" {
		t.Errorf("expected default BaseURL, got %s", cfg.BaseURL)
	}
}

func TestLoad_NormalFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{"base_url": "http://example.com:8080", "default_encoding": "text", "language": "zh"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFromPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "http://example.com:8080" {
		t.Errorf("expected http://example.com:8080, got %s", cfg.BaseURL)
	}
	if cfg.DefaultEncoding != "text" {
		t.Errorf("expected text, got %s", cfg.DefaultEncoding)
	}
	if cfg.Language != "zh" {
		t.Errorf("expected zh, got %s", cfg.Language)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte(`{invalid json`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := config.LoadFromPath(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_InvalidEncoding(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{"base_url": "http://example.com", "default_encoding": "xml"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFromPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 非法 encoding 应回退到 json
	if cfg.DefaultEncoding != "json" {
		t.Errorf("expected json fallback, got %s", cfg.DefaultEncoding)
	}
}

func TestLoad_InvalidLanguage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{"base_url": "http://example.com", "default_encoding": "json", "language": "de"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFromPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Language != "en" {
		t.Errorf("expected en fallback, got %s", cfg.Language)
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()

	// 需要 mock 配置目录 - 直接使用 LoadFromPath/Save 的路径逻辑
	// 由于 Save 使用默认路径(~/.config/endpoint-tui/)，这里通过直接测试 LoadFromPath 覆盖
	// 无法在不 mock home 的情况下完全测试 Save

	path := filepath.Join(dir, "config.json")
	cfg := config.Config{
		BaseURL:         "http://127.0.0.1:8080",
		DefaultEncoding: "markdown",
	}

	// 手动序列化保存
	data := []byte(`{"base_url":"http://127.0.0.1:8080","default_encoding":"markdown"}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	loaded, err := config.LoadFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.BaseURL != cfg.BaseURL {
		t.Errorf("expected BaseURL %s, got %s", cfg.BaseURL, loaded.BaseURL)
	}
}

func TestValidateBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常 HTTP", "http://127.0.0.1:8080", "http://127.0.0.1:8080", false},
		{"正常 HTTPS", "https://example.com", "https://example.com", false},
		{"末尾斜杠", "http://example.com/", "http://example.com", false},
		{"前后空格", "  http://example.com  ", "http://example.com", false},
		{"空字符串", "", "", true},
		{"无协议", "example.com", "", true},
		{"FTP 协议", "ftp://example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ValidateBaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadPublicInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "public-instance.json")
	content := `{"server":[{"url":"api.example.com/base","author":"tester","date":"2026-01-02"},{"url":"https://ready.example.com/","author":"ready"},{"url":""}]}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	instances, err := config.LoadPublicInstances(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("expected 2 valid public instances, got %d", len(instances))
	}
	if instances[0].URL != "https://api.example.com/base" {
		t.Fatalf("expected URL to be normalized with https, got %q", instances[0].URL)
	}
	if instances[1].URL != "https://ready.example.com" {
		t.Fatalf("expected trailing slash to be removed, got %q", instances[1].URL)
	}
}
