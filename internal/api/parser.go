// Package api handles endpoint discovery, parsing, and request execution.
package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Endpoint is the normalized endpoint representation.
type Endpoint struct {
	Name string
	Path string
}

// ParseEndpoints parses endpoint data in supported API response formats.
func ParseEndpoints(raw interface{}) ([]Endpoint, error) {
	if raw == nil {
		return nil, fmt.Errorf("endpoint field is missing or null")
	}

	switch v := raw.(type) {
	case []interface{}:
		return parseArray(v)
	case map[string]interface{}:
		return parseObject(v)
	default:
		return nil, fmt.Errorf("unsupported endpoint format: %T", raw)
	}
}

func parseArray(arr []interface{}) ([]Endpoint, error) {
	if len(arr) == 0 {
		return nil, fmt.Errorf("endpoint array is empty")
	}

	switch arr[0].(type) {
	case string:
		return parseStringArray(arr)
	case map[string]interface{}:
		return parseObjectArray(arr)
	default:
		return nil, fmt.Errorf("unsupported endpoint array element type: %T", arr[0])
	}
}

func parseStringArray(arr []interface{}) ([]Endpoint, error) {
	var result []Endpoint
	seen := make(map[string]bool)

	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		result = append(result, Endpoint{Name: s, Path: s})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint array has no valid endpoints")
	}
	return result, nil
}

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
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true

		if strings.TrimSpace(name) == "" {
			name = path
		}

		result = append(result, Endpoint{Name: name, Path: path})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint array has no valid endpoints")
	}
	return result, nil
}

func parseObject(obj map[string]interface{}) ([]Endpoint, error) {
	var result []Endpoint
	seen := make(map[string]bool)

	for name, val := range obj {
		path, ok := val.(string)
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true

		result = append(result, Endpoint{Name: name, Path: path})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("endpoint object has no valid endpoints")
	}
	return result, nil
}

// ParseJSONEndpoints extracts endpoint or endpoints fields from a JSON response.
func ParseJSONEndpoints(data []byte) ([]Endpoint, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	endpointData, ok := raw["endpoint"]
	if !ok {
		endpointData, ok = raw["endpoints"]
	}
	if !ok {
		return nil, fmt.Errorf("response is missing endpoint field")
	}

	return ParseEndpoints(endpointData)
}
