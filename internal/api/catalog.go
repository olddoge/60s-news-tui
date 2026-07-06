package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

//go:embed catalog.json
var endpointCatalogJSON []byte

// EndpointParam describes one GET query parameter for an endpoint.
type EndpointParam struct {
	Key      string            `json:"key"`
	Labels   map[string]string `json:"labels"`
	Required bool              `json:"required,omitempty"`
}

// EndpointMenuText holds localized display text and optional parameters.
type EndpointMenuText struct {
	Emoji  string            `json:"emoji"`
	Names  map[string]string `json:"names"`
	Params []EndpointParam   `json:"params,omitempty"`
}

// EndpointMenuTexts is the shared endpoint menu catalog.
var EndpointMenuTexts = mustLoadEndpointCatalog()

func mustLoadEndpointCatalog() map[string]EndpointMenuText {
	catalog, err := LoadEndpointCatalog(endpointCatalogJSON)
	if err != nil {
		panic(err)
	}
	return catalog
}

// LoadEndpointCatalog parses endpoint metadata from JSON.
func LoadEndpointCatalog(data []byte) (map[string]EndpointMenuText, error) {
	var catalog map[string]EndpointMenuText
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse endpoint catalog: %w", err)
	}
	return catalog, nil
}

// LocalizeEndpoints keeps only supported endpoints and replaces names with Chinese menu labels.
func LocalizeEndpoints(endpoints []Endpoint) []Endpoint {
	return LocalizeEndpointsForLanguage(endpoints, "zh")
}

// LocalizeEndpointsForLanguage keeps only supported endpoints and attaches localized menu labels and params.
func LocalizeEndpointsForLanguage(endpoints []Endpoint, language string) []Endpoint {
	localized := make([]Endpoint, 0, len(endpoints))
	seen := make(map[string]struct{}, len(endpoints))
	for _, ep := range endpoints {
		path := strings.TrimSpace(ep.Path)
		if path == "" {
			continue
		}

		menuText, ok := EndpointMenuTexts[lookupEndpointPath(path)]
		name := LocalizedEndpointName(menuText, language)
		if !ok || name == "" {
			continue
		}
		if menuText.Emoji != "" {
			name = strings.TrimSpace(menuText.Emoji + " " + name)
		}
		if _, exists := seen[path]; exists {
			continue
		}
		seen[path] = struct{}{}

		localized = append(localized, Endpoint{
			Name:   name,
			Path:   path,
			Params: normalizeEndpointParams(menuText.Params),
		})
	}
	return localized
}

// LocalizedEndpointName returns an endpoint menu name for language, falling back to Chinese.
func LocalizedEndpointName(menuText EndpointMenuText, language string) string {
	if language != "" {
		if name := strings.TrimSpace(menuText.Names[language]); name != "" {
			return name
		}
	}
	return strings.TrimSpace(menuText.Names["zh"])
}

// LocalizedParamLabel returns a parameter label for language, falling back to Chinese or key.
func LocalizedParamLabel(param EndpointParam, language string) string {
	if language != "" {
		if label := strings.TrimSpace(param.Labels[language]); label != "" {
			return label
		}
	}
	if label := strings.TrimSpace(param.Labels["zh"]); label != "" {
		return label
	}
	return param.Key
}

func normalizeEndpointParams(params []EndpointParam) []EndpointParam {
	normalized := make([]EndpointParam, 0, len(params))
	seen := make(map[string]struct{}, len(params))
	for _, param := range params {
		param.Key = strings.TrimSpace(param.Key)
		if param.Key == "" {
			continue
		}
		if _, exists := seen[param.Key]; exists {
			continue
		}
		seen[param.Key] = struct{}{}
		normalized = append(normalized, param)
	}
	return normalized
}

func lookupEndpointPath(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if u, err := url.Parse(endpoint); err == nil && u.Scheme != "" && u.Host != "" {
		return u.Path
	}
	if i := strings.IndexAny(endpoint, "?#"); i >= 0 {
		return endpoint[:i]
	}
	return endpoint
}
