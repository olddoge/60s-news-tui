package app

import "endpoint-tui/internal/api"

// Page is the current UI page.
type Page int

const (
	PageLoading Page = iota
	PageEndpointList
	PageEncodingSelect
	PageParamInput
	PageResult
	PageSettings
	PageError
)

// EndpointsLoadedMsg is sent when endpoint loading finishes.
type EndpointsLoadedMsg struct {
	Endpoints []api.Endpoint
	Error     error
}

// CurlResultMsg is sent when a curl request finishes.
type CurlResultMsg struct {
	Result api.CurlResult
}

// ErrorMsg carries a general error.
type ErrorMsg struct {
	Error error
}
