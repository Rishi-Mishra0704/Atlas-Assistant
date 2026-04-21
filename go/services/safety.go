package services

import "strings"

var dangerousApps = []string{"cmd", "powershell", "regedit", "taskmgr", "diskpart", "format"}

func IsSafeApp(name string) bool {
	n := strings.ToLower(name)
	for _, blocked := range dangerousApps {
		if strings.Contains(n, blocked) {
			return false
		}
	}
	return true
}
