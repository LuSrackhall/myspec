## ADDED Requirements

### Requirement: CLI entry point
The myspec binary SHALL accept a subcommand as the first argument and route to the corresponding handler.

#### Scenario: Valid subcommand
- **WHEN** user runs `myspec install /path/to/project`
- **THEN** the binary invokes the install handler with `/path/to/project` as argument

#### Scenario: No subcommand
- **WHEN** user runs `myspec` with no arguments
- **THEN** the binary prints help text listing all available subcommands

#### Scenario: Unknown subcommand
- **WHEN** user runs `myspec unknown-command`
- **THEN** the binary prints an error message and help text, exits with code 1

### Requirement: Help text
Each subcommand SHALL display usage information when invoked with `--help` or `-h`.

#### Scenario: Install help
- **WHEN** user runs `myspec install --help`
- **THEN** the binary prints usage: `myspec install <project-path>` with description

#### Scenario: Global help
- **WHEN** user runs `myspec --help`
- **THEN** the binary prints all subcommands with one-line descriptions
