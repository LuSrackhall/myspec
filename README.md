# myspec

Workflow management for Claude Code: define reusable development workflow skills and distribute them across projects via a Go CLI.

## What is myspec?

myspec solves a coordination problem in Claude Code projects. When you use Claude Code for development, the quality of output depends heavily on the process -- how requirements are explored, how designs are validated, and how implementation is verified. myspec packages these best practices into reusable, installable workflow skills.

myspec provides:

- **Structured brainstorming** (`myspec-br`) that turns vague ideas into approved designs through Socratic dialogue
- **Isolated workspaces** (`myspec-gwt`) that keep your main branch clean during experimentation
- **A custom OpenSpec schema** (`myspec-driven`) that chains brainstorming into the spec-driven development pipeline
- **A Go CLI** (planned) for installing and updating workflow skills across multiple projects

## Workflow Overview

A typical myspec-driven development cycle:

```
User has an idea
    |
    v
myspec-br (on main)
  - Explore project context
  - Ask clarifying questions (one at a time)
  - Propose 2-3 approaches with trade-offs
  - Present design section-by-section, get approval
  - Self-review, then user reviews
  - Write brainstorm-spec.md to OpenSpec change directory
    |
    v
myspec-br asks: "Create worktree?"
    |
    |-- Yes --> myspec-gwt creates worktree at .worktrees/change/<name>
    |           |
    |           v
    |       opsx:propose (in worktree)
    |         - Extract proposal from brainstorm-spec
    |         - Generate delta specs, design doc, task list
    |           |
    |           v
    |       opsx:apply (in worktree)
    |         - Work through tasks, implement changes
    |           |
    |           v
    |       opsx:verify (in worktree)
    |         - Validate implementation against specs
    |           |
    |           v
    |       Back on main:
    |         git merge change/<name>
    |           |
    |           v
    |       opsx:archive -> git commit -> git worktree remove
    |
    |-- No --> Stay on main (warning: may pollute main)
               |
               v
           opsx:propose -> apply -> verify -> archive (all on main)
```

**Discard path:** `git worktree remove .worktrees/change/<name> && git branch -d change/<name>`

## Skills

### myspec-br

**Location:** `.claude/skills/myspec-br/SKILL.md`

The brainstorming orchestrator. Transforms ideas into fully formed designs through natural collaborative dialogue. Key behaviors:

- HARD-GATE: no implementation until design is approved
- One question at a time, multiple choice preferred
- Always proposes 2-3 approaches before settling
- Section-by-section design approval with self-review
- Produces `brainstorm-spec.md` as the terminal artifact
- Decides whether to create a worktree, then delegates to myspec-gwt

### myspec-gwt

**Location:** `.claude/skills/myspec-gwt/SKILL.md`

Git worktree creation for change isolation. Handles:

- Detecting existing isolation (worktree or submodule)
- Validating the default branch (main or master)
- Checking for existing worktrees and repairing damaged ones
- Ensuring `.worktrees/` is gitignored
- Fallback to in-place work if worktree creation fails

## Custom Schema: myspec-driven

**Location:** `openspec/schemas/myspec-driven/`

A custom OpenSpec schema that integrates myspec-br into the spec-driven development pipeline.

### Artifact DAG

```
brainstorm-spec --> proposal --> specs + design --> tasks --> [apply] --> verify
```

| Artifact | Output | Produced By | Description |
|---|---|---|---|
| brainstorm-spec | `brainstorm-spec.md` | myspec-br | Approved design from Socratic dialogue |
| proposal | `proposal.md` | opsx:propose | Extracted Why/What/Capabilities/Impact |
| specs | `specs/**/*.md` | opsx:propose | Delta specifications (ADDED/MODIFIED/REMOVED) |
| design | `design.md` | opsx:propose | Technical implementation design |
| tasks | `tasks.md` | opsx:propose | Checkbox-based implementation checklist |
| verify | `verify.md` | opsx:verify | Post-implementation verification report |

The schema includes 6 templates in `openspec/schemas/myspec-driven/templates/`.

## Architecture

```
myspec/
|
|-- CLAUDE.md                 # Project constitution and constraints
|-- DESIGN.md                 # All design decisions (authoritative)
|-- README.md                 # This file
|
|-- .claude/
|   |-- skills/
|   |   |-- myspec-br/
|   |   |   `-- SKILL.md      # Brainstorming orchestrator skill
|   |   |-- myspec-gwt/
|   |   |   `-- SKILL.md      # Worktree creation skill
|   |   `-- openspec-*/        # OpenSpec workflow skills (installed by openspec init)
|   `-- commands/
|       `-- opsx/              # OpenSpec slash commands (installed by openspec init)
|
|-- openspec/
|   |-- schemas/
|   |   `-- myspec-driven/
|   |       |-- schema.yaml   # 6-artifact DAG definition
|   |       `-- templates/    # Template files for each artifact
|   |-- specs/                # Main branch authoritative specs (created at runtime)
|   `-- changes/              # Active change directories (created at runtime)
|
`-- .worktrees/               # Isolated git worktrees (gitignored, created at runtime)
```

## Installation

**Prerequisites:** Go 1.21+ installed ([download](https://go.dev/dl/))

```bash
# Install (recommended) — requires only Go, no clone needed
go install github.com/LuSrackhall/myspec@latest

# Verify installation
myspec --version
```

The binary is placed in `$GOPATH/bin` (or `$HOME/go/bin`). Ensure this directory is in your `PATH`.

### Build from source

```bash
git clone https://github.com/LuSrackhall/myspec.git
cd myspec
go build -o myspec .
```

### Usage

```bash
# Install skills to a project (auto-detects OpenSpec, auto-initializes if needed)
myspec install /path/to/project

# Update all installed projects
myspec update

# Update a specific project
myspec update /path/to/project

# List installed projects and versions
myspec list

# Uninstall from a project
myspec uninstall /path/to/project

# Check for available updates
myspec check

# Diagnose OpenSpec compatibility
myspec doctor
```

Installation copies skill files into the target project's `.claude/skills/` directory (no symlinks), making them git-trackable and self-contained.

The registry of installed projects is stored at `~/.config/myspec/registry.json`.

## Relationship to Other Tools

- **OpenSpec**: myspec uses OpenSpec for spec-driven change management and provides a custom schema. It does not modify any OpenSpec built-in schemas or skills.
- **Superpowers**: myspec-br borrows brainstorming methodology (Socratic dialogue, HARD-GATE, section-by-section approval) from the Superpowers project. No code or API dependency.

## Design Decisions

See [DESIGN.md](DESIGN.md) for all design decisions, including:

- Why brainstorming happens on main (no file output, no concurrency risk)
- Why implementation happens in a worktree (isolation, clean discard)
- Why myspec-br exists separately from opsx:explore (different scope and output)
- Why the Go CLI uses file copy over symlinks
- Finish workflow (merge then archive, not archive then merge)

## Status

Completed:
- myspec-br skill (brainstorming orchestrator)
- myspec-gwt skill (worktree creation)
- myspec-driven custom schema with 6 artifact DAG
- Schema template files (6 templates)

Planned:
- Go CLI development (install / update / list / uninstall / check)
- Global skill installation support
- Project-level override mechanism

## Git Conventions

All commit messages must be in **Chinese** using conventional commit format:

```
feat(simulation): 新增朝向系统组件
fix(render): 修复盾牌血条显示问题
docs: 更新设计文档
```

## License

TBD
