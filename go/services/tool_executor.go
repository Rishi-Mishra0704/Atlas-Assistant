package services

import "context"

type ToolExecutor interface {
	Execute(ctx context.Context, toolName string, input string) error
}

type toolExecutor struct {
	tools map[string]Tool
}

func NewToolExecutor(resolver AppResolver) ToolExecutor {
	executor := &toolExecutor{
		tools: make(map[string]Tool),
	}

	executor.register(&OpenAppTool{resolver: resolver})
	executor.register(&BrowserTool{})

	return executor
}

func (s *toolExecutor) register(tool Tool) {
	s.tools[tool.Name()] = tool
}

func (s *toolExecutor) Execute(ctx context.Context, toolName string, input string) error {
	tool, exists := s.tools[toolName]
	if !exists {
		return nil
	}

	return tool.Execute(ctx, input)
}
