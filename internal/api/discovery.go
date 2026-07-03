package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DefaultDiscoveryURL is used when ENDPOINT_DISCOVERY_URL is not set.
const DefaultDiscoveryURL = "http://localhost:13205"

// GetDiscoveryURL returns the endpoint discovery URL.
func GetDiscoveryURL() string {
	if url := os.Getenv("ENDPOINT_DISCOVERY_URL"); url != "" {
		return url
	}
	return DefaultDiscoveryURL
}

// FetchEndpoints gets the endpoint list from the discovery service.
func FetchEndpoints(discoveryURL string) ([]Endpoint, error) {
	if discoveryURL == "" {
		discoveryURL = GetDiscoveryURL()
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request endpoint list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("endpoint service returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read endpoint response: %w", err)
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("endpoint response is not valid JSON")
	}

	return ParseJSONEndpoints(body)
}
