package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DefaultDiscoveryURL 是默认的接口发现服务地址（可通过环境变量覆盖）。
const DefaultDiscoveryURL = "http://localhost:13205"

// GetDiscoveryURL 返回发现服务地址，优先使用环境变量。
func GetDiscoveryURL() string {
	if url := os.Getenv("ENDPOINT_DISCOVERY_URL"); url != "" {
		return url
	}
	return DefaultDiscoveryURL
}

// FetchEndpoints 从发现服务获取接口列表。
func FetchEndpoints(discoveryURL string) ([]Endpoint, error) {
	if discoveryURL == "" {
		discoveryURL = GetDiscoveryURL()
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("请求接口列表失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("接口列表服务返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取接口列表响应失败: %w", err)
	}

	// 尝试解析 JSON
	if !json.Valid(body) {
		return nil, fmt.Errorf("接口列表返回内容不是有效的 JSON")
	}

	return ParseJSONEndpoints(body)
}
