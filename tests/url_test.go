package tests

import (
	"net/url"
	"strings"
	"testing"

	"endpoint-tui/internal/urlutil"
)

// urlContainsParams 检查 URL 是否包含指定的查询参数。
func urlContainsParams(rawURL string, params map[string]string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	q := u.Query()
	for k, v := range params {
		if q.Get(k) != v {
			return false
		}
	}
	return true
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		endpoint   string
		encoding   string
		wantPrefix string // URL 前缀（不含查询参数顺序敏感部分）
		wantParams map[string]string
		wantErr    bool
	}{
		{
			name:       "正常拼接 - 根路径无末尾斜杠",
			baseURL:    "http://localhost:8080",
			endpoint:   "/api/test",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/test",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "根路径有末尾斜杠",
			baseURL:    "http://localhost:8080/",
			endpoint:   "/api/test",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/test",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "接口路径无前导斜杠",
			baseURL:    "http://localhost:8080",
			endpoint:   "api/test",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/test",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "根路径和接口路径都有斜杠",
			baseURL:    "http://localhost:8080/",
			endpoint:   "api/test",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/test",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "已有查询参数 - 正确追加 encoding",
			baseURL:    "http://localhost:8080",
			endpoint:   "/api/test?id=1",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/test",
			wantParams: map[string]string{"id": "1", "encoding": "json"},
		},
		{
			name:       "完整接口地址 - 不拼接根路径",
			baseURL:    "http://localhost:8080",
			endpoint:   "https://example.com/api/test",
			encoding:   "json",
			wantPrefix: "https://example.com/api/test",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "HTTPS 根路径",
			baseURL:    "https://example.com",
			endpoint:   "/api/secure",
			encoding:   "text",
			wantPrefix: "https://example.com/api/secure",
			wantParams: map[string]string{"encoding": "text"},
		},
		{
			name:       "encoding 为 markdown",
			baseURL:    "http://localhost:8080",
			endpoint:   "/api/docs",
			encoding:   "markdown",
			wantPrefix: "http://localhost:8080/api/docs",
			wantParams: map[string]string{"encoding": "markdown"},
		},
		{
			name:       "多级路径",
			baseURL:    "http://localhost:8080/api/v1",
			endpoint:   "/users/list",
			encoding:   "json",
			wantPrefix: "http://localhost:8080/api/v1/users/list",
			wantParams: map[string]string{"encoding": "json"},
		},
		{
			name:       "端口包含路径",
			baseURL:    "http://127.0.0.1:13205",
			endpoint:   "/api/article",
			encoding:   "json",
			wantPrefix: "http://127.0.0.1:13205/api/article",
			wantParams: map[string]string{"encoding": "json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := urlutil.BuildURL(tt.baseURL, tt.endpoint, tt.encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// 检查 URL 前缀
			if !strings.HasPrefix(got, tt.wantPrefix+"?") {
				t.Errorf("BuildURL() = %v, want prefix %v", got, tt.wantPrefix)
				return
			}

			// 检查查询参数
			if !urlContainsParams(got, tt.wantParams) {
				t.Errorf("BuildURL() = %v, missing expected params %v", got, tt.wantParams)
			}
		})
	}
}
