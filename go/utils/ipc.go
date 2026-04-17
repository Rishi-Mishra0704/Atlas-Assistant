package utils

import "net/http"

func NewSidecarClient() *http.Client {
	return http.DefaultClient
}
