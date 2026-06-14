## ADDED Requirements

### Requirement: Detect OpenSpec CLI
The system SHALL detect whether the OpenSpec CLI is available on the system PATH.

#### Scenario: OpenSpec found
- **WHEN** `openspec` is found on PATH via `exec.LookPath`
- **THEN** the system proceeds to version check

#### Scenario: OpenSpec not found
- **WHEN** `openspec` is not found on PATH
- **THEN** the system reports: "OpenSpec CLI not found. Install: npm install -g @fission-ai/openspec"
- **AND** returns an error

### Requirement: Version compatibility check
The system SHALL compare the installed OpenSpec version against the embedded `openspec-version.txt`.

#### Scenario: Versions match
- **WHEN** `openspec --version` output matches `openspec-version.txt` content
- **THEN** no warning is printed

#### Scenario: Installed version is newer
- **WHEN** installed OpenSpec version is newer than embedded version
- **THEN** the system prints a warning: "myspec was tested with OpenSpec <embedded>, you have <installed>"
- **AND** continues execution (does not block)

#### Scenario: Installed version is older
- **WHEN** installed OpenSpec version is older than embedded version
- **THEN** the system prints a warning with the exact npm install command to upgrade
- **AND** continues execution (does not block)

### Requirement: Auto-initialize OpenSpec in target project
The system SHALL detect whether the target project has been initialized with OpenSpec and initialize it if not.

#### Scenario: Project already initialized
- **WHEN** target project has `openspec/` directory with `config.yaml`
- **THEN** the system skips initialization

#### Scenario: Project not initialized
- **WHEN** target project does not have `openspec/` directory
- **THEN** the system runs `openspec init --tools claude` in the target project
- **AND** proceeds with installation after initialization succeeds
