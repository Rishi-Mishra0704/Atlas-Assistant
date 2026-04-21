package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BrowserTool struct {
	chromePath     string
	profileDir     string
	defaultProfile string
	profiles       map[string]string // display name → directory name
}

func NewBrowserTool(chromePath, profileDir, defaultProfile string) *BrowserTool {
	if chromePath == "" {
		chromePath = `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`
	}
	if profileDir == "" {
		profileDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data")
	}
	if defaultProfile == "" {
		defaultProfile = "Default"
	}

	b := &BrowserTool{
		chromePath:     chromePath,
		profileDir:     profileDir,
		defaultProfile: defaultProfile,
		profiles:       make(map[string]string),
	}
	b.scanProfiles()
	return b
}

func (b *BrowserTool) scanProfiles() {
	entries, err := os.ReadDir(b.profileDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name != "Default" && !strings.HasPrefix(name, "Profile") {
			continue
		}

		prefsPath := filepath.Join(b.profileDir, name, "Preferences")
		data, err := os.ReadFile(prefsPath)
		if err != nil {
			continue
		}

		var prefs struct {
			Profile struct {
				Name string `json:"name"`
			} `json:"profile"`
			AccountInfo []struct {
				Email string `json:"email"`
			} `json:"account_info"`
		}

		if err := json.Unmarshal(data, &prefs); err != nil {
			continue
		}

		// index by display name
		if prefs.Profile.Name != "" {
			b.profiles[strings.ToLower(prefs.Profile.Name)] = name
		}

		// index by email prefix (before @)
		if len(prefs.AccountInfo) > 0 && prefs.AccountInfo[0].Email != "" {
			email := prefs.AccountInfo[0].Email
			prefix := strings.Split(email, "@")[0]
			b.profiles[strings.ToLower(prefix)] = name
		}
	}
}

func (b *BrowserTool) resolveProfile(input string) (profile string, query string) {
	lower := strings.ToLower(input)

	for name, dir := range b.profiles {
		patterns := []string{
			"on my " + name + " profile",
			"on " + name + " profile",
			"in " + name + " profile",
			"using " + name,
			name + " profile",
		}
		for _, p := range patterns {
			if strings.Contains(lower, p) {
				cleaned := strings.Replace(lower, p, "", 1)
				return dir, strings.TrimSpace(cleaned)
			}
		}
	}

	return b.defaultProfile, input
}

func (b *BrowserTool) Name() string {
	return "open_browser"
}

func (b *BrowserTool) Execute(ctx context.Context, input string) error {
	profile, query := b.resolveProfile(input)

	args := []string{
		"--profile-directory=" + profile,
		"--user-data-dir=" + b.profileDir,
	}

	if strings.HasPrefix(query, "http") || strings.Contains(query, ".com") || strings.Contains(query, ".org") || strings.Contains(query, ".io") {
		args = append(args, "--new-tab", query)
	} else if query != "" {
		searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", strings.ReplaceAll(query, " ", "+"))
		args = append(args, "--new-tab", searchURL)
	}
	cmd := exec.Command(b.chromePath, args...)
	return cmd.Start()
}
