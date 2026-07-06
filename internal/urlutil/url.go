// Package urlutil builds safe endpoint URLs.
package urlutil

import (
	"net/url"
	"strings"
)

// BuildURL combines a base URL, endpoint path, and encoding query parameter.
func BuildURL(baseURL, endpoint, encoding string) (string, error) {
	return BuildURLWithParams(baseURL, endpoint, encoding, nil)
}

// BuildURLWithParams combines a base URL, endpoint path, extra GET params, and encoding query parameter.
func BuildURLWithParams(baseURL, endpoint, encoding string, params map[string]string) (string, error) {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return appendQuery(endpoint, encoding, params)
	}

	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	full := baseURL + endpoint
	return appendQuery(full, encoding, params)
}

func appendQuery(rawURL, encoding string, params map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range params {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		q.Set(key, value)
	}
	q.Set("encoding", encoding)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
