package services

import (
	"atlas/config"
	"context"
	"strings"
)

type ToolExecutor interface {
	Execute(ctx context.Context, toolName string, input string) error
}

type toolExecutor struct {
	tools map[string]Tool
}

func NewToolExecutor(resolver AppResolver, cfg config.Config) ToolExecutor {
	executor := &toolExecutor{
		tools: make(map[string]Tool),
	}

	executor.register(&OpenAppTool{resolver: resolver})
	executor.register(NewBrowserTool(cfg.ChromePath, cfg.ChromeProfileDir, cfg.ChromeDefaultProfile))

	return executor
}
func (s *toolExecutor) register(tool Tool) {
	s.tools[tool.Name()] = tool
}

func (s *toolExecutor) Execute(ctx context.Context, toolName string, input string) error {
	// catch chrome/browser misroutes
	lower := strings.ToLower(input)
	if toolName == "open_app" && (strings.Contains(lower, "chrome") || strings.Contains(lower, "browser")) {
		toolName = "open_browser"
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return nil
	}

	return tool.Execute(ctx, input)
}
