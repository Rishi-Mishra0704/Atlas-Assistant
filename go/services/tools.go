package services

import (
	"context"
	"fmt"
	"os/exec"
)

type Tool interface {
	Name() string
	Execute(ctx context.Context, input string) error
}

type OpenAppTool struct {
	resolver AppResolver
}

func (t *OpenAppTool) Name() string {
	return "open_app"
}

func (t *OpenAppTool) Execute(ctx context.Context, input string) error {
	if !IsSafeApp(input) {
		return fmt.Errorf("blocked: %s is not allowed", input)
	}

	resolved, err := t.resolver.Resolve(input)
	if err != nil {
		// No shortcut found — try as a direct command
		cmd := exec.Command("cmd", "/c", "start", "", input)
		return cmd.Start()
	}

	return t.resolver.Open(resolved)
}

type BrowserTool struct{}

func (t *BrowserTool) Name() string {
	return "open_browser"
}

func (t *BrowserTool) Execute(ctx context.Context, input string) error {
	cmd := exec.Command("cmd", "/c", "start", "https://google.com")
	return cmd.Start()
}
