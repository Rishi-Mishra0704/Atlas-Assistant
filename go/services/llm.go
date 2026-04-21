package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// ---- types ----

type IntentDecision struct {
	Route      string  `json:"route"`
	Intent     string  `json:"intent"`
	Closure    bool    `json:"closure"`
	Confidence float64 `json:"confidence,omitempty"`
}

const (
	IntentRouteChat  = "chat"
	IntentRouteTool  = "tool"
	IntentRouteClose = "close"
)

// ---- types ----

type ToolCommand struct {
	Action string `json:"action"`
	Target string `json:"target"`
}

// ---- interface ----

type ILLMService interface {
	DecideIntent(ctx context.Context, text string) (IntentDecision, error)
	GenerateReply(ctx context.Context, userText string, decision IntentDecision) (string, error)
	GenerateCommand(ctx context.Context, userText string, decision IntentDecision) (ToolCommand, error)
}

// ---- constructor ----

func NewLLMService(baseURL, model string) ILLMService {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://localhost:11434"
	}
	if strings.TrimSpace(model) == "" {
		model = "qwen2.5:0.5b"
	}

	return &llmService{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ---- implementation ----

type llmService struct {
	baseURL string
	model   string
	client  *http.Client
}

//
// ---------------- TIER 1 ----------------
//

func (c *llmService) DecideIntent(ctx context.Context, text string) (IntentDecision, error) {

	system := `
You are Atlas' Tier 1 intent router.

Your ONLY job is to classify intent.

Return JSON:

{
  "route": "chat | tool | close",
  "intent": "short_snake_case_label",
  "tool": "optional_tool_name",
  "closure": true_or_false
}

Rules:
- Do NOT generate replies
- Only classify intent

Routing:
- chat: conversation/questions
- tool: user wants an action
- close: user explicitly ends session

Close examples:
- bye
- goodbye
- exit
- stop

IMPORTANT:
- closure = true ONLY if route = close
- otherwise closure = false

Return ONLY JSON
`

	req := map[string]any{
		"model":      c.model,
		"system":     system,
		"prompt":     text,
		"stream":     false,
		"format":     "json",
		"keep_alive": "10m",
	}

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return IntentDecision{}, err
	}
	defer resp.Body.Close()

	var ollamaResp struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return IntentDecision{}, err
	}

	var decision IntentDecision
	if err := json.Unmarshal([]byte(ollamaResp.Response), &decision); err != nil {
		// fallback
		return IntentDecision{
			Route:   IntentRouteChat,
			Intent:  "fallback",
			Closure: false,
		}, nil
	}

	decision = normalizeDecision(decision, text)

	return decision, nil
}

//
// ---------------- TIER 2 ----------------
//

func (c *llmService) GenerateReply(
	ctx context.Context,
	userText string,
	decision IntentDecision,
) (string, error) {

	system := `
You are Atlas, a local AI assistant.

You will be given:
- user_text
- route (chat | tool | close)
- intent

Rules:
- Always produce a short, natural reply (<= 15 words).
- If answering a factual question, you may extend to 30–40 words.
- If route == "tool":
  - Acknowledge briefly (e.g., "Opening VS Code.").
- If route == "close":
  - Say a brief goodbye.
- If route == "chat":
  - Answer naturally and briefly.
- Do NOT mention routes or intents.
- Do NOT output JSON, only plain text.
`

	payload := map[string]any{
		"user_text": userText,
		"route":     decision.Route,
		"intent":    decision.Intent,
	}

	input, _ := json.Marshal(payload)

	req := map[string]any{
		"model":      "llama3.1:8b",
		"system":     system,
		"prompt":     string(input),
		"stream":     false,
		"keep_alive": "10m",
	}

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ollamaResp struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

//
// ---------------- TIER 2b: COMMAND EXTRACTION ----------------
//

func (c *llmService) GenerateCommand(ctx context.Context, userText string, decision IntentDecision) (ToolCommand, error) {
	system := `You are Atlas's command extractor.

Given the user text and intent, return a structured JSON command.
Only describe what to do, never output shell commands.

Return JSON:
{
  "action": "open_app",
  "target": "the app name the user wants"
}

Rules:
- action is ALWAYS "open_app". No exceptions. Never use open_url, open_discord, launch_game, or anything else.
- target is the app name the user said, nothing more
- Return ONLY valid JSON

`

	payload := map[string]any{
		"user_text": userText,
		"intent":    decision.Intent,
	}
	input, _ := json.Marshal(payload)

	req := map[string]any{
		"model":      "llama3.1:8b",
		"system":     system,
		"prompt":     string(input),
		"stream":     false,
		"format":     "json",
		"keep_alive": "10m",
	}

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return ToolCommand{}, err
	}
	defer resp.Body.Close()

	var ollamaResp struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return ToolCommand{}, err
	}

	var cmd ToolCommand
	if err := json.Unmarshal([]byte(ollamaResp.Response), &cmd); err != nil {
		return ToolCommand{}, err
	}

	return cmd, nil
}

// ---------------- NORMALIZATION ----------------
func normalizeDecision(decision IntentDecision, text string) IntentDecision {
	decision.Route = strings.ToLower(strings.TrimSpace(decision.Route))
	decision.Intent = strings.ToLower(strings.TrimSpace(decision.Intent))

	// ---- enforce route from intent ----
	if isToolIntent(decision.Intent, text) {
		decision.Route = IntentRouteTool
	}

	if isCloseIntent(decision.Intent, text) {
		decision.Route = IntentRouteClose
		decision.Closure = true
	}

	// ---- enforce closure ----
	if decision.Route != IntentRouteClose {
		decision.Closure = false
	} else {
		decision.Closure = true
	}

	return decision
}

func isExplicitClose(text string) bool {
	t := strings.ToLower(text)

	return strings.Contains(t, "bye") ||
		strings.Contains(t, "goodbye") ||
		strings.Contains(t, "exit") ||
		strings.Contains(t, "stop")
}

func isToolIntent(intent string, text string) bool {
	intent = strings.ToLower(intent)
	text = strings.ToLower(text)

	return strings.Contains(intent, "open") ||
		strings.Contains(intent, "launch") ||
		strings.Contains(intent, "run") ||
		strings.Contains(text, "open") ||
		strings.Contains(text, "launch") ||
		strings.Contains(text, "run")
}

func isCloseIntent(intent string, text string) bool {
	intent = strings.ToLower(intent)
	text = strings.ToLower(text)

	return strings.Contains(intent, "goodbye") ||
		strings.Contains(intent, "exit") ||
		strings.Contains(intent, "stop") ||
		strings.Contains(text, "bye") ||
		strings.Contains(text, "exit") ||
		strings.Contains(text, "stop")
}
