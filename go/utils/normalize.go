package utils

import "strings"

func DetectToolFromText(text string) string {
	t := strings.ToLower(text)

	switch {
	case strings.Contains(t, "vs code"),
		strings.Contains(t, "vscode"),
		strings.Contains(t, "code editor"),
		strings.Contains(t, "editor"):
		return "open_vscode"

	case strings.Contains(t, "chrome"),
		strings.Contains(t, "browser"):
		return "open_browser"
	}

	return ""
}
