package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AppResolver interface {
	Resolve(target string) (string, error)
	Open(resolved string) error
}

type appResolver struct {
	index map[string]string // lowercase name → full .lnk path
}

func NewAppResolver() AppResolver {
	r := &appResolver{index: make(map[string]string)}
	r.buildIndex()
	return r
}

func (r *appResolver) buildIndex() {
	dirs := []string{
		filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
		`C:\ProgramData\Microsoft\Windows\Start Menu\Programs`,
	}

	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if strings.EqualFold(filepath.Ext(path), ".lnk") {
				name := strings.ToLower(strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())))
				r.index[name] = path
			}
			return nil
		})
	}

	// common aliases
	aliases := map[string]string{
		"vscode":  "visual studio code",
		"vs code": "visual studio code",
		"chrome":  "google chrome",
	}
	for alias, real := range aliases {
		if path, ok := r.index[real]; ok {
			r.index[alias] = path
		}
	}
}

func (r *appResolver) Resolve(target string) (string, error) {
	t := strings.ToLower(strings.TrimSpace(target))

	if path, ok := r.index[t]; ok {
		return path, nil
	}

	for name, path := range r.index {
		if strings.Contains(name, t) || strings.Contains(t, name) {
			return path, nil
		}
	}

	return "", fmt.Errorf("app not found: %s", target)
}

func (r *appResolver) Open(resolved string) error {
	cmd := exec.Command("cmd", "/c", "start", "", resolved)
	return cmd.Start()
}
