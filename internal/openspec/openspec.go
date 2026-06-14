package openspec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Detect checks if OpenSpec CLI is available and returns its version.
func Detect() (found bool, version string, err error) {
	path, err := exec.LookPath("openspec")
	if err != nil {
		return false, "", nil
	}
	_ = path

	out, err := exec.Command("openspec", "--version").Output()
	if err != nil {
		return true, "", fmt.Errorf("openspec found but version check failed: %w", err)
	}

	return true, strings.TrimSpace(string(out)), nil
}

// CheckVersion compares installed OpenSpec version against the expected version.
// Returns: match (true if exact match), warning message (if mismatch).
func CheckVersion(installedVersion, expectedVersion string) (match bool, warning string) {
	if installedVersion == expectedVersion {
		return true, ""
	}

	if installedVersion > expectedVersion {
		return false, fmt.Sprintf(
			"Warning: myspec was tested with OpenSpec %s, but you have %s.\n"+
				"Skills may work correctly. If you experience issues, run:\n"+
				"  npm install -g @fission-ai/openspec@%s",
			expectedVersion, installedVersion, expectedVersion,
		)
	}

	return false, fmt.Sprintf(
		"Warning: myspec was tested with OpenSpec %s, but you have %s (older).\n"+
			"Please upgrade:\n"+
			"  npm install -g @fission-ai/openspec@%s",
		expectedVersion, installedVersion, expectedVersion,
	)
}

// InitProject runs `openspec init --tools claude` in the given directory.
func InitProject(projectPath string) error {
	cmd := exec.Command("openspec", "init", "--tools", "claude")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInitialized checks if a project has been initialized with OpenSpec.
func IsInitialized(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "openspec", "config.yaml"))
	return err == nil
}

// SetSchema updates the schema field in openspec/config.yaml.
// If the file doesn't exist, creates it with just the schema field.
// Preserves existing context and rules fields.
func SetSchema(projectPath, schemaName string) error {
	configPath := filepath.Join(projectPath, "openspec", "config.yaml")

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// Create new config
		return os.WriteFile(configPath, []byte("schema: "+schemaName+"\n"), 0644)
	}
	if err != nil {
		return fmt.Errorf("cannot read config.yaml: %w", err)
	}

	content := string(data)

	// Replace existing schema line
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "schema:") {
			lines[i] = "schema: " + schemaName
			found = true
			break
		}
	}

	if !found {
		// Prepend schema field
		lines = append([]string{"schema: " + schemaName, ""}, lines...)
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
}
