package utils

type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}
