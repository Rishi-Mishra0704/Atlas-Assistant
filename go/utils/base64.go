package utils

import (
	"encoding/base64"
	"fmt"
)

func EncodeTokenToBase64(token string) string {
	return base64.URLEncoding.EncodeToString([]byte(token))
}

func DecodeTokenFromBase64(encoded string) (string, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 token: %w", err)
	}
	return string(decodedBytes), nil
}
