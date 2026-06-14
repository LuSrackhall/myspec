# cli-install Specification

## Purpose
TBD - created by archiving change build-go-cli. Update Purpose after archive.
## Requirements
### Requirement: Install skills to project
The `myspec install` command SHALL copy embedded skill files and schema to the target project directory.

#### Scenario: Fresh install
- **WHEN** user runs `myspec install /path/to/project`
- **THEN** the following files are copied:
  - `embed/skills/myspec-br/SKILL.md` → `/path/to/project/.claude/skills/myspec-br/SKILL.md`
  - `embed/skills/myspec-gwt/SKILL.md` → `/path/to/project/.claude/skills/myspec-gwt/SKILL.md`
  - `embed/schemas/myspec-driven/` → `/path/to/project/openspec/schemas/myspec-driven/`
- **AND** `openspec/config.yaml` is created or updated with `schema: myspec-driven`
- **AND** the install is recorded in `~/.config/myspec/registry.json`

#### Scenario: Target project not initialized with OpenSpec
- **WHEN** user runs `myspec install /path/to/project` and the project has no `openspec/` directory
- **THEN** the command automatically runs `openspec init --tools claude` in the target project
- **THEN** proceeds with the normal install

#### Scenario: OpenSpec CLI not installed
- **WHEN** user runs `myspec install` and `openspec` is not found on PATH
- **THEN** the command prints: "OpenSpec CLI not found. Install it: npm install -g @fission-ai/openspec"
- **AND** exits with code 1

#### Scenario: Already installed
- **WHEN** user runs `myspec install /path/to/project` and the project is already in the registry
- **THEN** the command overwrites existing files and updates the registry entry

### Requirement: Update skills in project
The `myspec update` command SHALL replace all skill files and schema in the target project with the current embedded versions.

#### Scenario: Update specific project
- **WHEN** user runs `myspec update /path/to/project`
- **THEN** all embedded files are replaced (full replacement)
- **AND** the registry version is updated

#### Scenario: Update all projects
- **WHEN** user runs `myspec update` with no arguments
- **THEN** all projects in the registry are updated sequentially

### Requirement: Uninstall skills from project
The `myspec uninstall` command SHALL remove all myspec-installed files from the target project.

#### Scenario: Uninstall
- **WHEN** user runs `myspec uninstall /path/to/project`
- **THEN** `.claude/skills/myspec-br/` and `.claude/skills/myspec-gwt/` are removed
- **AND** `openspec/schemas/myspec-driven/` is removed
- **AND** the entry is removed from registry.json

