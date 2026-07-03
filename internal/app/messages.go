package app

import "endpoint-tui/internal/api"

// Page 表示当前页面状态。
type Page int

const (
	PageLoading        Page = iota // 加载接口列表
	PageEndpointList               // 接口列表页
	PageEncodingSelect             // encoding 选择页
	PageResult                     // 请求结果页
	PageSettings                   // 设置页
	PageError                      // 错误页
)

// EndpointsLoadedMsg 接口列表加载完成消息。
type EndpointsLoadedMsg struct {
	Endpoints []api.Endpoint
	Error     error
}

// CurlResultMsg curl 请求完成消息。
type CurlResultMsg struct {
	Result api.CurlResult
}

// ErrorMsg 通用错误消息。
type ErrorMsg struct {
	Error error
}
