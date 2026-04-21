package services

import "atlas/config"

type Service struct {
	LLM  ILLMService
	Tool ToolExecutor
}

func NewService(cfg config.Config) *Service {
	resolver := NewAppResolver()
	return &Service{
		LLM:  NewLLMService(cfg.OllamaURL, cfg.OllamaIntentModel),
		Tool: NewToolExecutor(resolver),
	}
}
