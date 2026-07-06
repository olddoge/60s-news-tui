package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// PublicInstance describes one server entry from public-instance.json.
type PublicInstance struct {
	URL    string `json:"url"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

type publicInstanceFile struct {
	Servers []PublicInstance `json:"server"`
}

var embeddedPublicInstances []byte

// SetDefaultPublicInstances registers embedded public-instance.json data.
func SetDefaultPublicInstances(data []byte) {
	embeddedPublicInstances = data
}

// LoadPublicInstances reads public server entries from public-instance.json.
func LoadPublicInstances(path string) ([]PublicInstance, error) {
	var data []byte
	if path == "" {
		path = findPublicInstancePath()
	}
	if path != "" {
		var err error
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read public instances: %w", err)
		}
	} else if len(embeddedPublicInstances) > 0 {
		data = embeddedPublicInstances
	} else {
		return nil, nil
	}

	var file publicInstanceFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("invalid public instances JSON: %w", err)
	}

	instances := make([]PublicInstance, 0, len(file.Servers))
	seen := make(map[string]struct{}, len(file.Servers))
	for _, instance := range file.Servers {
		instance.URL = NormalizePublicInstanceURL(instance.URL)
		instance.Author = strings.TrimSpace(instance.Author)
		instance.Date = strings.TrimSpace(instance.Date)
		if instance.URL == "" {
			continue
		}
		if _, exists := seen[instance.URL]; exists {
			continue
		}
		seen[instance.URL] = struct{}{}
		instances = append(instances, instance)
	}
	return instances, nil
}

// NormalizePublicInstanceURL turns a public instance host/path into a base URL.
func NormalizePublicInstanceURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/")
}

func findPublicInstancePath() string {
	const filename = "public-instance.json"
	if path := findUpward(filename); path != "" {
		return path
	}
	if exe, err := os.Executable(); err == nil {
		path := filepath.Join(filepath.Dir(exe), filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func findUpward(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for i := 0; i < 6; i++ {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
