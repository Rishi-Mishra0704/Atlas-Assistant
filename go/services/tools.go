package services

import (
	"context"
	"os/exec"
)

type Tool interface {
	Name() string
	Execute(ctx context.Context, input string) error
}

type VSCodeTool struct{}

func (t *VSCodeTool) Name() string {
	return "open_vscode"
}

func (t *VSCodeTool) Execute(ctx context.Context, input string) error {
	cmd := exec.Command("cmd", "/c", "start", "", "code")
	return cmd.Start()
}

type BrowserTool struct{}

func (t *BrowserTool) Name() string {
	return "open_browser"
}

func (t *BrowserTool) Execute(ctx context.Context, input string) error {
	cmd := exec.Command("cmd", "/c", "start", "https://google.com")
	return cmd.Start()
}
