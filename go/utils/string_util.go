package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

func FromText(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
