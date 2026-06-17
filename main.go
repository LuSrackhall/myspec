package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LuSrackhall/myspec/internal/emb"
	"github.com/LuSrackhall/myspec/internal/openspec"
	"github.com/LuSrackhall/myspec/internal/registry"
)

//go:embed all:embed
var embedFS embed.FS

//go:embed openspec-version.txt
var expectedOpenSpecVersion string

const version = "0.4.0"

var commands = map[string]struct {
	run  func(args []string)
	help string
}{
	"install":   {run: cmdInstall, help: "Install skills to a project: myspec install <project-path>"},
	"update":    {run: cmdUpdate, help: "Update installed skills: myspec update [project-path]"},
	"list":      {run: cmdList, help: "List installed projects: myspec list"},
	"uninstall": {run: cmdUninstall, help: "Uninstall skills: myspec uninstall <project-path>"},
	"check":     {run: cmdCheck, help: "Check for outdated versions: myspec check"},
	"doctor":    {run: cmdDoctor, help: "Diagnose OpenSpec compatibility: myspec doctor"},
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	name := os.Args[1]

	if name == "--help" || name == "-h" {
		printUsage()
		return
	}

	if name == "--version" || name == "-v" {
		fmt.Printf("myspec %s\n", version)
		return
	}

	cmd, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", name)
		printUsage()
		os.Exit(1)
	}

	args := os.Args[2:]
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Printf("Usage: myspec %s\n\n%s\n", name, cmd.help)
			return
		}
	}

	cmd.run(args)
}

func printUsage() {
	fmt.Println("myspec - Claude Code workflow skill manager")
	fmt.Println()
	fmt.Println("Usage: myspec <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	for name, cmd := range commands {
		fmt.Printf("  %-12s %s\n", name, cmd.help)
	}
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Printf("  %-12s %s\n", "--help, -h", "Show help")
	fmt.Printf("  %-12s %s\n", "--version, -v", "Show version")
}

func cmdInstall(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: myspec install <project-path>")
		os.Exit(1)
	}

	projectPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid path: %v\n", err)
		os.Exit(1)
	}

	// 1. Check OpenSpec CLI
	found, ver, err := openspec.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking OpenSpec: %v\n", err)
		os.Exit(1)
	}
	if !found {
		fmt.Fprintln(os.Stderr, "OpenSpec CLI not found. Install it:")
		fmt.Fprintln(os.Stderr, "  npm install -g @fission-ai/openspec")
		os.Exit(1)
	}

	// 2. Version check
	expected := strings.TrimSpace(expectedOpenSpecVersion)
	if match, warning := openspec.CheckVersion(ver, expected); !match {
		fmt.Fprintln(os.Stderr, warning)
		fmt.Fprintln(os.Stderr)
	}

	// 3. Auto-init OpenSpec if needed
	if !openspec.IsInitialized(projectPath) {
		fmt.Println("Initializing OpenSpec in target project...")
		if err := openspec.InitProject(projectPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: OpenSpec init failed: %v\n", err)
			os.Exit(1)
		}
	}

	// 4. Set schema to myspec-driven
	if err := openspec.SetSchema(projectPath, "myspec-driven"); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting schema: %v\n", err)
		os.Exit(1)
	}

	// 5. Copy skills
	if err := emb.CopySkills(embedFS, projectPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying skills: %v\n", err)
		os.Exit(1)
	}

	// 5. Copy schema
	if err := emb.CopySchema(embedFS, projectPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying schema: %v\n", err)
		os.Exit(1)
	}

	// 6. Update registry
	r, err := registry.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	r.Set(projectPath, registry.Entry{
		Version:     version,
		InstalledAt: time.Now().UTC(),
		Skills:      []string{"myspec-br", "myspec-gwt"},
		Schema:      "myspec-driven",
	})

	if err := r.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving registry: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installed myspec skills to %s\n", projectPath)
	fmt.Printf("  Skills: myspec-br, myspec-gwt\n")
	fmt.Printf("  Schema: myspec-driven\n")
	fmt.Printf("  Version: %s\n", version)
}

func cmdUpdate(args []string) {
	r, err := registry.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	paths := []string{}
	if len(args) > 0 {
		p, _ := filepath.Abs(args[0])
		paths = []string{p}
	} else {
		for p := range r.Installed {
			paths = append(paths, p)
		}
	}

	if len(paths) == 0 {
		fmt.Println("No projects to update.")
		return
	}

	for _, p := range paths {
		if err := emb.CopySkills(embedFS, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", p, err)
			continue
		}
		if err := emb.CopySchema(embedFS, p); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating schema in %s: %v\n", p, err)
			continue
		}
		r.Set(p, registry.Entry{
			Version:     version,
			InstalledAt: time.Now().UTC(),
			Skills:      []string{"myspec-br", "myspec-gwt"},
			Schema:      "myspec-driven",
		})
		fmt.Printf("Updated: %s\n", p)
	}

	r.Save()
}

func cmdList(args []string) {
	r, err := registry.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	if len(r.Installed) == 0 {
		fmt.Println("No projects installed.")
		return
	}

	for p, e := range r.Installed {
		fmt.Printf("  %-50s %s\n", p, e.Version)
	}
}

func cmdUninstall(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: myspec uninstall <project-path>")
		os.Exit(1)
	}

	projectPath, _ := filepath.Abs(args[0])

	r, err := registry.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	if _, ok := r.Get(projectPath); !ok {
		fmt.Fprintf(os.Stderr, "Project not found in registry: %s\n", projectPath)
		os.Exit(1)
	}

	emb.RemoveSkills(projectPath)
	emb.RemoveSchema(projectPath)
	r.Remove(projectPath)
	r.Save()

	fmt.Printf("Uninstalled: %s\n", projectPath)
}

func cmdCheck(args []string) {
	r, err := registry.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	if len(r.Installed) == 0 {
		fmt.Println("No projects installed.")
		return
	}

	for p, e := range r.Installed {
		if e.Version == version {
			fmt.Printf("  %-50s %s (up to date)\n", p, e.Version)
		} else {
			fmt.Printf("  %-50s %s (current: %s)\n", p, e.Version, version)
		}
	}
}

func cmdDoctor(args []string) {
	fmt.Println("myspec doctor")
	fmt.Println("=============")

	expected := strings.TrimSpace(expectedOpenSpecVersion)

	// Check OpenSpec
	found, ver, err := openspec.Detect()
	if !found {
		fmt.Println("  [FAIL] OpenSpec CLI: not found")
		fmt.Println("         Install: npm install -g @fission-ai/openspec")
	} else if err != nil {
		fmt.Printf("  [WARN] OpenSpec CLI: found but error: %v\n", err)
	} else {
		match, warning := openspec.CheckVersion(ver, expected)
		if match {
			fmt.Printf("  [OK]   OpenSpec CLI: %s\n", ver)
		} else {
			fmt.Printf("  [WARN] OpenSpec CLI: %s (expected %s)\n", ver, expected)
			fmt.Println("         " + warning)
		}
	}

	// Check registry
	r, err := registry.Load()
	if err != nil {
		fmt.Printf("  [FAIL] Registry: %v\n", err)
	} else {
		fmt.Printf("  [OK]   Registry: %d project(s)\n", len(r.Installed))
	}
}
