## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init github.com/LuSrackhalllu/myspec`)
- [x] 1.2 Create directory structure (`cmd/myspec/`, `internal/`, `embed/`)
- [x] 1.3 Create `openspec-version.txt` with current OpenSpec version
- [x] 1.4 Copy skill files and schema to `embed/` directory
- [x] 1.5 Add `.gitignore` entries for Go build artifacts

## 2. Core CLI Framework (cli-core)

- [x] 2.1 Implement main entry point with subcommand routing (`cmd/myspec/main.go`)
- [x] 2.2 Implement help text for global and per-subcommand usage
- [x] 2.3 Implement `--help` / `-h` flag handling for each subcommand

## 3. Registry Management (cli-registry)

- [x] 3.1 Implement registry file read/write (`internal/registry/registry.go`)
- [x] 3.2 Implement `myspec list` command
- [x] 3.3 Implement `myspec check` command (version comparison)
- [x] 3.4 Implement `myspec doctor` command (OpenSpec diagnostics)

## 4. OpenSpec Bridge (cli-openspec-bridge)

- [x] 4.1 Implement OpenSpec CLI detection (`exec.LookPath`)
- [x] 4.2 Implement version comparison against `openspec-version.txt`
- [x] 4.3 Implement `openspec init --tools claude` auto-execution
- [x] 4.4 Implement config.yaml merge logic (update schema field, preserve context/rules)

## 5. Install / Update / Uninstall (cli-install)

- [x] 5.1 Implement `go:embed` resource loading (`internal/embed/embed.go`)
- [x] 5.2 Implement file copy logic (skills + schema to target project)
- [x] 5.3 Implement `myspec install` (full flow: detect OpenSpec → init → copy → registry)
- [x] 5.4 Implement `myspec update` (single project + all projects)
- [x] 5.5 Implement `myspec uninstall` (remove files + registry entry)

## 6. Testing

- [ ] 6.1 Unit tests for registry read/write
- [ ] 6.2 Unit tests for version comparison
- [ ] 6.3 Unit tests for config.yaml merge
- [ ] 6.4 Integration test: install to temp directory, verify files
- [ ] 6.5 Integration test: update replaces files correctly
- [ ] 6.6 Integration test: uninstall removes all files

## 7. Build and Distribution

- [x] 7.1 Verify `go build -o myspec .` produces working binary
- [x] 7.2 Verify binary works on macOS (primary platform)
- [ ] 7.3 Document installation in README.md
