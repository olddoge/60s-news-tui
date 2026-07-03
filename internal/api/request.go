package api

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// CurlResult 表示一次 curl 请求的结果。
type CurlResult struct {
	URL        string
	StatusCode int
	Stdout     string
	Stderr     string
	ExitCode   int
	Duration   time.Duration
	Error      error
	Cancelled  bool
}

// CommandExecutor 定义命令执行接口，便于单元测试 mock。
type CommandExecutor interface {
	Execute(ctx context.Context, url string) CurlResult
}

// CurlExecutor 通过系统 curl 命令执行 HTTP 请求。
type CurlExecutor struct {
	Timeout       time.Duration
	ConnectTimout time.Duration
}

// NewCurlExecutor 创建默认的 CurlExecutor。
func NewCurlExecutor() *CurlExecutor {
	return &CurlExecutor{
		Timeout:       60 * time.Second,
		ConnectTimout: 10 * time.Second,
	}
}

// CheckCurlAvailable 检查系统是否安装了 curl。
func CheckCurlAvailable() error {
	_, err := exec.LookPath("curl")
	if err != nil {
		return fmt.Errorf("未找到 curl 命令。\n\n请在 Debian 中执行：\nsudo apt update && sudo apt install -y curl")
	}
	return nil
}

// Execute 执行 curl 请求，支持通过 context 取消。
func (e *CurlExecutor) Execute(ctx context.Context, requestURL string) CurlResult {
	result := CurlResult{
		URL: requestURL,
	}

	start := time.Now()

	cmd := exec.CommandContext(
		ctx,
		"curl",
		"--silent",
		"--show-error",
		"--location",
		"--connect-timeout",
		fmt.Sprintf("%.0f", e.ConnectTimout.Seconds()),
		"--max-time",
		fmt.Sprintf("%.0f", e.Timeout.Seconds()),
		requestURL,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Duration = time.Since(start)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if err != nil {
		// 判断是否被取消
		if ctx.Err() != nil {
			result.Cancelled = true
			result.Error = fmt.Errorf("请求已取消")
			return result
		}

		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Error = err
		return result
	}

	result.ExitCode = 0
	return result
}
