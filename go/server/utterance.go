package server

import (
	"atlas/utils"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UtteranceRequest struct {
	Text string `json:"text"`
}

type UtteranceResponse struct {
	Reply      string  `json:"reply"`
	Closure    bool    `json:"closure"`
	Route      string  `json:"route"`
	Intent     string  `json:"intent"`
	Tool       string  `json:"tool,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

func (s *Server) HandleSpeech(ctx echo.Context) error {
	var req UtteranceRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err, "Invalid request body", 400))
	}

	// ---- Tier 1: Intent Routing ----
	decision, err := s.services.LLM.DecideIntent(ctx.Request().Context(), req.Text)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err, "Failed to route utterance with Ollama", http.StatusBadGateway))
	}

	// ---- Tool Execution (deterministic) ----
	if decision.Route == "tool" {
		tool := utils.DetectToolFromText(req.Text)

		if tool == "" {
			log.Printf("No tool detected from text: %s", req.Text)
			decision.Route = "chat"
		} else {
			log.Printf("Executing tool: %s", tool)
			_ = s.services.Tool.Execute(ctx.Request().Context(), tool, req.Text)
		}
	}

	// ---- Tier 2: Reply Generation ----
	reply, err := s.services.LLM.GenerateReply(ctx.Request().Context(), req.Text, decision)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.ErrorResponse(err, "Failed to generate reply", http.StatusBadGateway))
	}

	resp := UtteranceResponse{
		Reply:      reply,
		Closure:    decision.Closure,
		Route:      decision.Route,
		Intent:     decision.Intent,
		Confidence: decision.Confidence,
	}

	return ctx.JSON(http.StatusOK, utils.SuccessResponse(resp))
}
