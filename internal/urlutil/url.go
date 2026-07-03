// Package urlutil 提供安全的 URL 拼接和查询参数处理。
package urlutil

import (
	"net/url"
	"strings"
)

// BuildURL 将根路径和接口路径组合为完整 URL，并追加 encoding 查询参数。
// 如果 endpoint 已经是完整 HTTP/HTTPS 地址，则直接使用该地址，
// 不再拼接根路径。
func BuildURL(baseURL, endpoint, encoding string) (string, error) {
	// 如果 endpoint 已经是完整地址，直接使用
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return appendEncoding(endpoint, encoding)
	}

	// 清理根路径
	baseURL = strings.TrimRight(baseURL, "/")

	// 确保 endpoint 以 "/" 开头
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	full := baseURL + endpoint
	return appendEncoding(full, encoding)
}

// appendEncoding 向 URL 追加 encoding 查询参数，正确处理已有查询参数的情况。
func appendEncoding(rawURL, encoding string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("encoding", encoding)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
