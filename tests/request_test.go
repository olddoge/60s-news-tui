package tests

import (
	"context"
	"testing"
	"time"

	"endpoint-tui/internal/api"
)

// MockExecutor 是用于测试的 CommandExecutor mock。
type MockExecutor struct {
	Result api.CurlResult
}

func (m *MockExecutor) Execute(ctx context.Context, url string) api.CurlResult {
	return m.Result
}

func TestMockExecutor_Success(t *testing.T) {
	mock := &MockExecutor{
		Result: api.CurlResult{
			URL:       "http://example.com/api/test?encoding=json",
			Stdout:    `{"code":200,"message":"success"}`,
			Stderr:    "",
			ExitCode:  0,
			Duration:  100 * time.Millisecond,
			Error:     nil,
			Cancelled: false,
		},
	}

	result := mock.Execute(context.Background(), "http://example.com/api/test?encoding=json")
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if result.Stdout != `{"code":200,"message":"success"}` {
		t.Errorf("unexpected stdout: %s", result.Stdout)
	}
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
}

func TestMockExecutor_Error(t *testing.T) {
	mock := &MockExecutor{
		Result: api.CurlResult{
			URL:      "http://example.com/notfound",
			Stdout:   "",
			Stderr:   "curl: (7) Failed to connect",
			ExitCode: 7,
			Error:    &mockExitError{exitCode: 7},
		},
	}

	result := mock.Execute(context.Background(), "http://example.com/notfound")
	if result.ExitCode != 7 {
		t.Errorf("expected exit code 7, got %d", result.ExitCode)
	}
}

func TestMockExecutor_Cancelled(t *testing.T) {
	mock := &MockExecutor{
		Result: api.CurlResult{
			URL:       "http://example.com/slow",
			Cancelled: true,
			Error:     context.Canceled,
		},
	}

	result := mock.Execute(context.Background(), "http://example.com/slow")
	if !result.Cancelled {
		t.Error("expected cancelled to be true")
	}
}

func TestCurlExecutor_NewCurlExecutor(t *testing.T) {
	executor := api.NewCurlExecutor()
	if executor.Timeout != 60*time.Second {
		t.Errorf("expected 60s timeout, got %v", executor.Timeout)
	}
	if executor.ConnectTimout != 10*time.Second {
		t.Errorf("expected 10s connect timeout, got %v", executor.ConnectTimout)
	}
}

func TestCheckCurlAvailable(t *testing.T) {
	err := api.CheckCurlAvailable()
	if err != nil {
		// curl 未安装在测试环境中是正常的
		t.Logf("curl not available (may be expected): %v", err)
	}
}

// mockExitError 模拟 exec.ExitError。
type mockExitError struct {
	exitCode int
}

func (e *mockExitError) Error() string {
	return "exit status " + itoa(e.exitCode)
}

func (e *mockExitError) ExitCode() int {
	return e.exitCode
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
