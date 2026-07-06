package tests

import (
	"testing"

	"endpoint-tui/internal/api"
)

func TestLocalizeEndpoints_UsesChineseNamesAndFiltersUnsupported(t *testing.T) {
	endpoints := []api.Endpoint{
		{Name: "/v2/60s", Path: "/v2/60s"},
		{Name: "/v2/chemical", Path: "/v2/chemical"},
		{Name: "/v2/qrcode", Path: "/v2/qrcode?size=256"},
		{Name: "/unknown", Path: "/unknown"},
		{Name: "remote", Path: "https://example.com/v2/bili"},
	}

	got := api.LocalizeEndpoints(endpoints)
	if len(got) != 3 {
		t.Fatalf("expected 3 localized endpoints, got %d: %#v", len(got), got)
	}

	want := []struct {
		name string
		path string
	}{
		{name: "📰 每天 60 秒读懂世界", path: "/v2/60s"},
		{name: "🔳 生成二维码", path: "/v2/qrcode?size=256"},
		{name: "📺 哔哩哔哩热搜", path: "https://example.com/v2/bili"},
	}
	for i, item := range want {
		if got[i].Name != item.name || got[i].Path != item.path {
			t.Fatalf("endpoint %d mismatch: got %#v, want name=%q path=%q", i, got[i], item.name, item.path)
		}
	}
}

func TestEndpointMenuTexts_IsPublicCatalog(t *testing.T) {
	text, ok := api.EndpointMenuTexts["/v2/60s"]
	if !ok {
		t.Fatal("expected /v2/60s in public endpoint menu catalog")
	}
	if text.Names["zh"] != "每天 60 秒读懂世界" || text.Names["en"] != "Understand the World in 60 Seconds" || text.Emoji != "📰" {
		t.Fatalf("unexpected public catalog text: %#v", text)
	}
}

func TestEndpointMenuTexts_LoadsParamsFromJSON(t *testing.T) {
	text, ok := api.EndpointMenuTexts["/v2/qrcode"]
	if !ok {
		t.Fatal("expected /v2/qrcode in public endpoint menu catalog")
	}
	if len(text.Params) != 3 {
		t.Fatalf("expected qrcode params from JSON, got %d", len(text.Params))
	}
	param := text.Params[0]
	if param.Key != "text" || param.Required {
		t.Fatalf("unexpected qrcode text param metadata: %#v", param)
	}
	if api.LocalizedParamLabel(param, "en") != "QR Code Text" {
		t.Fatalf("expected English param label, got %q", api.LocalizedParamLabel(param, "en"))
	}
	if api.LocalizedParamLabel(param, "zh") != "二维码内容" {
		t.Fatalf("expected Chinese param label, got %q", api.LocalizedParamLabel(param, "zh"))
	}
}

func TestLocalizeEndpointsForLanguage_UsesEnglishNames(t *testing.T) {
	got := api.LocalizeEndpointsForLanguage([]api.Endpoint{
		{Name: "/v2/60s", Path: "/v2/60s"},
		{Name: "/v2/answer", Path: "/v2/answer"},
	}, "en")
	if len(got) != 2 {
		t.Fatalf("expected 2 localized endpoints, got %d", len(got))
	}
	if got[0].Name != "📰 Understand the World in 60 Seconds" {
		t.Fatalf("expected English 60s name, got %q", got[0].Name)
	}
	if got[1].Name != "📖 Random Book of Answers" {
		t.Fatalf("expected English answer name, got %q", got[1].Name)
	}
}

func TestLocalizeEndpoints_AllowsEmptyEmoji(t *testing.T) {
	original := api.EndpointMenuTexts["/v2/60s"]
	api.EndpointMenuTexts["/v2/60s"] = api.EndpointMenuText{Names: original.Names}
	defer func() {
		api.EndpointMenuTexts["/v2/60s"] = original
	}()

	got := api.LocalizeEndpoints([]api.Endpoint{{Name: "/v2/60s", Path: "/v2/60s"}})
	if len(got) != 1 {
		t.Fatalf("expected endpoint to remain visible without emoji, got %d", len(got))
	}
	if got[0].Name != "每天 60 秒读懂世界" {
		t.Fatalf("expected plain Chinese name without emoji, got %q", got[0].Name)
	}
}
