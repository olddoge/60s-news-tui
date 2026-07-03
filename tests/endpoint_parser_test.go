package tests

import (
	"encoding/json"
	"testing"

	"endpoint-tui/internal/api"
)

func TestParseEndpoints_StringArray(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{"/api/article", "/api/document", "/api/status"},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 3 {
		t.Fatalf("expected 3 endpoints, got %d", len(endpoints))
	}
	if endpoints[0].Path != "/api/article" {
		t.Errorf("expected /api/article, got %s", endpoints[0].Path)
	}
	if endpoints[0].Name != "/api/article" {
		t.Errorf("expected name to be path, got %s", endpoints[0].Name)
	}
}

func TestParseEndpoints_ObjectArray(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{
			map[string]interface{}{"name": "获取文章", "path": "/api/article"},
			map[string]interface{}{"name": "获取文档", "path": "/api/document"},
		},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
	if endpoints[0].Name != "获取文章" {
		t.Errorf("expected '获取文章', got '%s'", endpoints[0].Name)
	}
	if endpoints[0].Path != "/api/article" {
		t.Errorf("expected /api/article, got %s", endpoints[0].Path)
	}
}

func TestParseEndpoints_KeyValueObject(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": map[string]interface{}{
			"article":  "/api/article",
			"document": "/api/document",
		},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
}

func TestParseEndpoints_Empty(t *testing.T) {
	// 空数组
	data := map[string]interface{}{
		"endpoint": []interface{}{},
	}
	raw, _ := extractEndpoint(data)
	_, err := api.ParseEndpoints(raw)
	if err == nil {
		t.Error("expected error for empty array")
	}
}

func TestParseEndpoints_MissingEndpoint(t *testing.T) {
	_, err := api.ParseJSONEndpoints([]byte(`{"other": "data"}`))
	if err == nil {
		t.Error("expected error for missing endpoint")
	}
}

func TestParseEndpoints_DuplicateFilter(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{"/api/article", "/api/article", "/api/document"},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 unique endpoints, got %d", len(endpoints))
	}
}

func TestParseEndpoints_EmptyStringFilter(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{"/api/article", "", "  "},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 valid endpoint, got %d", len(endpoints))
	}
}

func TestParseEndpoints_InvalidType(t *testing.T) {
	_, err := api.ParseEndpoints(42)
	if err == nil {
		t.Error("expected error for invalid type")
	}
	_, err = api.ParseEndpoints("string")
	if err == nil {
		t.Error("expected error for string type")
	}
}

func TestParseEndpoints_WrongArrayElementType(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{123, 456},
	}
	raw, _ := extractEndpoint(data)
	_, err := api.ParseEndpoints(raw)
	if err == nil {
		t.Error("expected error for non-string array elements")
	}
}

func TestParseJSONEndpoints_MalformedJSON(t *testing.T) {
	_, err := api.ParseJSONEndpoints([]byte(`{this is not json`))
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestParseEndpoints_ObjectArrayMissingName(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": []interface{}{
			map[string]interface{}{"path": "/api/test"},
		},
	}
	raw, err := extractEndpoint(data)
	if err != nil {
		t.Fatal(err)
	}
	endpoints, err := api.ParseEndpoints(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	// 没有 name 时应使用 path 作为 name
	if endpoints[0].Name != "/api/test" {
		t.Errorf("expected name fallback to path, got '%s'", endpoints[0].Name)
	}
}

// extractEndpoint 从 map 中提取 endpoint 字段。
func extractEndpoint(data map[string]interface{}) (interface{}, error) {
	if _, ok := data["endpoint"]; !ok {
		return nil, nil
	}
	// 需要支持深层解析场景，把数据转成 JSON 再解析回来以模拟真实解析过程
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result["endpoint"], nil
}
