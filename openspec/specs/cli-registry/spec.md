# cli-registry Specification

## Purpose
TBD - created by archiving change build-go-cli. Update Purpose after archive.
## Requirements
### Requirement: Registry file format
The registry SHALL be stored at `~/.config/myspec/registry.json` as a JSON file.

#### Scenario: Registry structure
- **WHEN** the registry file is read
- **THEN** it contains a `version` field (integer) and an `installed` object
- **AND** each key in `installed` is an absolute project path
- **AND** each value contains `version` (string), `installedAt` (ISO 8601), `skills` (string array), and `schema` (string)

#### Scenario: Registry file does not exist
- **WHEN** `myspec list` runs and `~/.config/myspec/registry.json` does not exist
- **THEN** the command prints "No projects installed" and exits cleanly

### Requirement: List installed projects
The `myspec list` command SHALL display all installed projects with their versions.

#### Scenario: List with installed projects
- **WHEN** user runs `myspec list`
- **THEN** the command prints each project path and version in a tabular format

#### Scenario: List with no projects
- **WHEN** user runs `myspec list` and no projects are installed
- **THEN** the command prints "No projects installed"

### Requirement: Check for outdated versions
The `myspec check` command SHALL compare installed versions against the current git tag.

#### Scenario: All up to date
- **WHEN** user runs `myspec check` and all installed versions match the current tag
- **THEN** the command prints "All projects up to date"

#### Scenario: Some outdated
- **WHEN** user runs `myspec check` and some projects have older versions
- **THEN** the command prints each outdated project with its installed version and current version

### Requirement: Doctor diagnostics
The `myspec doctor` command SHALL check system compatibility.

#### Scenario: All checks pass
- **WHEN** user runs `myspec doctor`
- **THEN** the command checks: OpenSpec CLI exists, OpenSpec version matches, registry is readable
- **AND** prints a diagnostic report with pass/fail for each check

