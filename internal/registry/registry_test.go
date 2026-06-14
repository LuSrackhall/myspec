package registry

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMissingFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	r, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if r.Version != 1 {
		t.Errorf("expected version 1, got %d", r.Version)
	}
	if len(r.Installed) != 0 {
		t.Errorf("expected 0 installed, got %d", len(r.Installed))
	}
}

func TestSaveAndLoad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	r := &Registry{Version: 1, Installed: map[string]Entry{}}
	r.Set("/test/project", Entry{
		Version:     "v1.0.0",
		InstalledAt: time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
		Skills:      []string{"myspec-br", "myspec-gwt"},
		Schema:      "myspec-driven",
	})

	if err := r.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	r2, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	entry, ok := r2.Get("/test/project")
	if !ok {
		t.Fatal("expected entry for /test/project")
	}
	if entry.Version != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %s", entry.Version)
	}
	if entry.Schema != "myspec-driven" {
		t.Errorf("expected schema myspec-driven, got %s", entry.Schema)
	}
	if len(entry.Skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(entry.Skills))
	}
}

func TestRemove(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	r := &Registry{Version: 1, Installed: map[string]Entry{}}
	r.Set("/a", Entry{Version: "v1.0.0"})
	r.Set("/b", Entry{Version: "v1.0.0"})
	r.Save()

	r.Remove("/a")
	r.Save()

	r2, _ := Load()
	if _, ok := r2.Get("/a"); ok {
		t.Error("expected /a to be removed")
	}
	if _, ok := r2.Get("/b"); !ok {
		t.Error("expected /b to still exist")
	}
}

func TestCorruptJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".config", "myspec")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "registry.json"), []byte("not json"), 0644)

	_, err := Load()
	if err == nil {
		t.Error("expected error for corrupt JSON")
	}
}

func TestGetMissing(t *testing.T) {
	r := &Registry{Version: 1, Installed: map[string]Entry{}}
	_, ok := r.Get("/nonexistent")
	if ok {
		t.Error("expected false for missing entry")
	}
}
