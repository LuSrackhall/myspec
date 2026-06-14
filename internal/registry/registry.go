package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installedAt"`
	Skills      []string  `json:"skills"`
	Schema      string    `json:"schema"`
}

type Registry struct {
	Version   int              `json:"version"`
	Installed map[string]Entry `json:"installed"`
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "myspec", "registry.json"), nil
}

func Load() (*Registry, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &Registry{Version: 1, Installed: map[string]Entry{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read registry: %w", err)
	}

	var r Registry
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("cannot parse registry: %w", err)
	}
	if r.Installed == nil {
		r.Installed = map[string]Entry{}
	}
	return &r, nil
}

func (r *Registry) Save() error {
	p, err := path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot serialize registry: %w", err)
	}

	return os.WriteFile(p, data, 0644)
}

func (r *Registry) Set(projectPath string, entry Entry) {
	r.Installed[projectPath] = entry
}

func (r *Registry) Remove(projectPath string) {
	delete(r.Installed, projectPath)
}

func (r *Registry) Get(projectPath string) (Entry, bool) {
	e, ok := r.Installed[projectPath]
	return e, ok
}
