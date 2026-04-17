package services

import "atlas/config"

type Service struct {
	Orchestrator *OrchestratorService
	Voice        *VoiceService
	Sidecars     *SidecarService
	Tools        *ToolsService
	LLM          *LLMService
	RAG          *RAGService
	Memory       *MemoryService
	Permissions  *PermissionsService
	Audit        *AuditService
}

func NewService(cfg config.Config) *Service {
	return &Service{
		Orchestrator: &OrchestratorService{},
		Voice:        &VoiceService{},
		Sidecars:     &SidecarService{},
		Tools:        &ToolsService{},
		LLM:          &LLMService{},
		RAG:          &RAGService{QdrantURL: cfg.QdrantURL},
		Memory:       &MemoryService{},
		Permissions:  &PermissionsService{},
		Audit:        &AuditService{},
	}
}
