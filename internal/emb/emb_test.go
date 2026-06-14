package emb

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

//go:embed embed
var testFS embed.FS

func TestCopySkills(t *testing.T) {
	dir := t.TempDir()
	if err := CopySkills(testFS, dir); err != nil {
		t.Fatalf("CopySkills() failed: %v", err)
	}

	for _, name := range []string{"myspec-br", "myspec-gwt"} {
		path := filepath.Join(dir, ".claude", "skills", name, "SKILL.md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file at %s", path)
		}
	}
}

func TestCopySchema(t *testing.T) {
	dir := t.TempDir()
	if err := CopySchema(testFS, dir); err != nil {
		t.Fatalf("CopySchema() failed: %v", err)
	}

	schemaPath := filepath.Join(dir, "openspec", "schemas", "myspec-driven", "schema.yaml")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Errorf("expected schema at %s", schemaPath)
	}
}

func TestRemoveSkills(t *testing.T) {
	dir := t.TempDir()
	CopySkills(testFS, dir)

	if err := RemoveSkills(dir); err != nil {
		t.Fatalf("RemoveSkills() failed: %v", err)
	}

	for _, name := range []string{"myspec-br", "myspec-gwt"} {
		path := filepath.Join(dir, ".claude", "skills", name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("expected %s to be removed", path)
		}
	}
}

func TestRemoveSchema(t *testing.T) {
	dir := t.TempDir()
	CopySchema(testFS, dir)

	if err := RemoveSchema(dir); err != nil {
		t.Fatalf("RemoveSchema() failed: %v", err)
	}

	path := filepath.Join(dir, "openspec", "schemas", "myspec-driven")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed", path)
	}
}

func TestCopySkillsToNonexistentDir(t *testing.T) {
	dir := t.TempDir()
	// Target .claude dir doesn't exist yet - should create it
	if err := CopySkills(testFS, dir); err != nil {
		t.Fatalf("CopySkills() should create dirs: %v", err)
	}
}

// Test helpers - these verify the testdata directory structure matches what we expect
func TestTestdataExists(t *testing.T) {
	entries, err := fs.ReadDir(testFS, "embed/skills")
	if err != nil {
		t.Fatalf("embed/skills not found: %v", err)
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 skill dirs, got %d", len(entries))
	}
}
