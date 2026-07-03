package api

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// CurlResult represents one curl request result.
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

// CommandExecutor executes requests and can be mocked in tests.
type CommandExecutor interface {
	Execute(ctx context.Context, url string) CurlResult
}

// CurlExecutor executes HTTP requests through the system curl command.
type CurlExecutor struct {
	Timeout       time.Duration
	ConnectTimout time.Duration
}

// NewCurlExecutor creates the default curl executor.
func NewCurlExecutor() *CurlExecutor {
	return &CurlExecutor{
		Timeout:       60 * time.Second,
		ConnectTimout: 10 * time.Second,
	}
}

// CheckCurlAvailable verifies that curl is installed.
func CheckCurlAvailable() error {
	_, err := exec.LookPath("curl")
	if err != nil {
		return fmt.Errorf("curl command not found\n\nInstall curl first, for example on Debian:\nsudo apt update && sudo apt install -y curl")
	}
	return nil
}

// Execute runs a curl request and supports cancellation through context.
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
		if ctx.Err() != nil {
			result.Cancelled = true
			result.Error = fmt.Errorf("request cancelled")
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
