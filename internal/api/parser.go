// Package api 负责接口列表的发现、解析以及 curl 请求的执行。
package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Endpoint 表示一个接口，内部统一使用此结构。
type Endpoint struct {
	Name string
	Path string
}

// ParseEndpoints 解析 endpoint 字段，兼容三种格式：
// 1. 字符串数组: ["/api/a", "/api/b"]
// 2. 对象数组: [{"name":"A","path":"/api/a"}, ...]
// 3. 键值对象: {"a":"/api/a", "b":"/api/b"}
func ParseEndpoints(raw interface{}) ([]Endpoint, error) {
	if raw == nil {
		return nil, fmt.Errorf("endpoint 字段不存在或为 null")
	}

	switch v := raw.(type) {
	case []interface{}:
		return parseArray(v)
	case map[string]interface{}:
		return parseObject(v)
	default:
		return nil, fmt.Errorf("不支持的 endpoint 格式: %T", raw)
	}
}

// parseArray 处理 endpoint 为数组的情况。
func parseArray(arr []interface{}) ([]Endpoint, error) {
	if len(arr) == 0 {
		return nil, fmt.Errorf("endpoint 数组为空")
	}

	// 判断第一个元素的类型
	switch arr[0].(type) {
	case string:
		return parseStringArray(arr)
	case map[string]interface{}:
		return parseObjectArray(arr)
	default:
		return nil, fmt.Errorf("不支持的 endpoint 数组元素类型: %T", arr[0])
	}
}

// parseStringArray 处理 ["/api/a", "/api/b"] 格式。
func parseStringArray(arr []interface{}) ([]Endpoint, error) {
	var result []Endpoint
	seen := make(map[string]bool)

	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if seen[s] {
			continue
		}
		seen[s] = true
		result = append(result, Endpoint{Name: s, Path: s})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint 数组中没有有效接口")
	}
	return result, nil
}

// parseObjectArray 处理 [{"name":"A","path":"/api/a"}, ...] 格式。
func parseObjectArray(arr []interface{}) ([]Endpoint, error) {
	var result []Endpoint
	seen := make(map[string]bool)

	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := obj["name"].(string)
		path, _ := obj["path"].(string)

		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		if seen[path] {
			continue
		}
		seen[path] = true

		if strings.TrimSpace(name) == "" {
			name = path
		}

		result = append(result, Endpoint{Name: name, Path: path})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint 数组中没有有效接口")
	}
	return result, nil
}

// parseObject 处理 {"a":"/api/a", "b":"/api/b"} 格式。
func parseObject(obj map[string]interface{}) ([]Endpoint, error) {
	var result []Endpoint
	seen := make(map[string]bool)

	for name, val := range obj {
		path, ok := val.(string)
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		if seen[path] {
			continue
		}
		seen[path] = true

		result = append(result, Endpoint{Name: name, Path: path})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint 对象中没有有效接口")
	}
	return result, nil
}

// ParseJSONEndpoints 从 JSON 字节中提取并解析 endpoint 字段。
// 兼容 endpoint（单数）和 endpoints（复数）两种字段名。
func ParseJSONEndpoints(data []byte) ([]Endpoint, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	// 尝试两种字段名
	endpointData, ok := raw["endpoint"]
	if !ok {
		endpointData, ok = raw["endpoints"]
	}
	if !ok {
		return nil, fmt.Errorf("返回数据中缺少 endpoint 字段")
	}

	return ParseEndpoints(endpointData)
}
