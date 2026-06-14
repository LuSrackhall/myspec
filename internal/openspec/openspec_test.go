package openspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckVersionEqual(t *testing.T) {
	match, warning := CheckVersion("1.4.1", "1.4.1")
	if !match {
		t.Error("expected match=true for equal versions")
	}
	if warning != "" {
		t.Errorf("expected empty warning, got %q", warning)
	}
}

func TestCheckVersionNewer(t *testing.T) {
	match, warning := CheckVersion("1.5.0", "1.4.1")
	if match {
		t.Error("expected match=false for newer version")
	}
	if warning == "" {
		t.Error("expected warning for newer version")
	}
}

func TestCheckVersionOlder(t *testing.T) {
	match, warning := CheckVersion("1.3.0", "1.4.1")
	if match {
		t.Error("expected match=false for older version")
	}
	if warning == "" {
		t.Error("expected warning for older version")
	}
}

func TestCheckVersionMultiDigit(t *testing.T) {
	match, _ := CheckVersion("1.10.0", "1.4.1")
	if match {
		t.Error("expected match=false")
	}
	// 1.10.0 > 1.4.1 (10 > 4), should be "newer"
	_, warning := CheckVersion("1.10.0", "1.4.1")
	if warning == "" {
		t.Error("expected warning for 1.10.0 vs 1.4.1")
	}
}

func TestCheckVersionOlderMultiDigit(t *testing.T) {
	// 1.4.1 < 1.10.0, should be "older"
	match, warning := CheckVersion("1.4.1", "1.10.0")
	if match {
		t.Error("expected match=false")
	}
	if warning == "" {
		t.Error("expected warning")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.10.0", "1.4.1", 1},  // key: 10 > 4
		{"1.4.1", "1.10.0", -1}, // key: 4 < 10
		{"2.0.0", "1.9.9", 1},
		{"0.9.0", "0.10.0", -1},
	}
	for _, tt := range tests {
		got := compareVersions(tt.a, tt.b)
		if (got > 0 && tt.want <= 0) || (got < 0 && tt.want >= 0) || (got == 0 && tt.want != 0) {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSetSchemaNoFile(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "openspec"), 0755)

	if err := SetSchema(dir, "myspec-driven"); err != nil {
		t.Fatalf("SetSchema() failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	if string(data) != "schema: myspec-driven\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestSetSchemaExistingValue(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "openspec"), 0755)
	os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte("schema: spec-driven\n"), 0644)

	if err := SetSchema(dir, "myspec-driven"); err != nil {
		t.Fatalf("SetSchema() failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	expected := "schema: myspec-driven\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestSetSchemaPreservesContext(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "openspec"), 0755)
	content := "schema: spec-driven\ncontext: |\n  Tech stack: Go\n"
	os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(content), 0644)

	SetSchema(dir, "myspec-driven")

	data, _ := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	s := string(data)
	if !contains(s, "schema: myspec-driven") {
		t.Error("missing schema line")
	}
	if !contains(s, "context:") {
		t.Error("context field was lost")
	}
}

func TestSetSchemaSkipsComment(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "openspec"), 0755)
	content := "# schema: old-value\nschema: spec-driven\n"
	os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(content), 0644)

	SetSchema(dir, "myspec-driven")

	data, _ := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	s := string(data)
	if !contains(s, "# schema: old-value") {
		t.Error("comment was modified")
	}
	if !contains(s, "schema: myspec-driven") {
		t.Error("schema not updated")
	}
}

func TestSetSchemaNoSchemaLine(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "openspec"), 0755)
	content := "context: |\n  Tech stack: Go\n"
	os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(content), 0644)

	SetSchema(dir, "myspec-driven")

	data, _ := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	s := string(data)
	if !contains(s, "schema: myspec-driven") {
		t.Error("schema line not prepended")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
