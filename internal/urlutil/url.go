// Package urlutil builds safe endpoint URLs.
package urlutil

import (
	"net/url"
	"strings"
)

// BuildURL combines a base URL, endpoint path, and encoding query parameter.
func BuildURL(baseURL, endpoint, encoding string) (string, error) {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return appendEncoding(endpoint, encoding)
	}

	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	full := baseURL + endpoint
	return appendEncoding(full, encoding)
}

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
