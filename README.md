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

### myspec-apply

**Location:** `.claude/skills/myspec-apply/SKILL.md`

Wraps OpenSpec's apply workflow with automatic git commits per task group. Handles:

- Implementing tasks by group with automatic conventional commits
- Using the user's preferred language for commit messages
- Prompting for verification after all tasks are complete

### myspec-verify

**Location:** `.claude/skills/myspec-verify/SKILL.md`

Wraps OpenSpec's verify workflow with user acceptance and iteration. Handles:

- Document verification (Completeness/Correctness/Coherence)
- User acceptance checkpoint (must explicitly confirm)
- Iteration decision loop when user is not satisfied
- Backfilling all artifacts to match final implementation

### myspec-merge

**Location:** `.claude/skills/myspec-merge/SKILL.md`

Handles the complete post-verification merge workflow. Handles:

- Main branch sync check (local vs origin, user decides)
- Post-sync re-verification via myspec-verify
- Merge method selection (merge commit / squash / rebase)
- Archive and worktree cleanup

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

### Download binary (recommended)

Download the latest binary for your platform from [GitHub Releases](https://github.com/LuSrackhall/myspec/releases):

| Platform | File |
|----------|------|
| macOS (Intel) | `myspec-darwin-amd64` |
| macOS (Apple Silicon) | `myspec-darwin-arm64` |
| Linux (x86_64) | `myspec-linux-amd64` |
| Linux (ARM64) | `myspec-linux-arm64` |
| Windows (x86_64) | `myspec-windows-amd64.exe` |
| Windows (ARM64) | `myspec-windows-arm64.exe` |

```bash
# macOS / Linux: rename and add to PATH
mv myspec-<os>-<arch> myspec
chmod +x myspec
sudo mv myspec /usr/local/bin/

# Windows: move to a directory in your PATH
# e.g., C:\Users\<you>\bin\myspec.exe
```

### Install via Go

**Prerequisites:** Go 1.21+ installed ([download](https://go.dev/dl/))

```bash
go install github.com/LuSrackhall/myspec@latest
```

**Note for users in China:** If `go install` times out, configure the Go proxy:

```bash
go env -w GOPROXY=https://goproxy.cn,direct GONOSUMDB=*
go install github.com/LuSrackhall/myspec@latest
```

The binary is placed in `$GOPATH/bin` (or `$HOME/go/bin`).

**If `myspec` command is not found**, add Go's bin directory to your PATH:

```bash
# zsh (macOS default)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc

# bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc && source ~/.bashrc

# Windows (PowerShell)
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$(go env GOPATH)\bin", "User")
```

### Windows Note

myspec-gwt (worktree creation) uses bash shell commands. On Windows, use **Git Bash**, **WSL**, or **MSYS2**. The Go CLI itself runs natively on Windows.

### Quick Start

```bash
# 1. Install myspec
go install github.com/LuSrackhall/myspec@latest

# 2. Install skills to your project
myspec install /path/to/your/project

# 3. Verify everything works
myspec doctor
```

`myspec install` automatically detects OpenSpec CLI and initializes it if needed. It copies skill files (`myspec-br`, `myspec-gwt`, `myspec-apply`, `myspec-verify`, `myspec-merge`) and the `myspec-driven` schema into your project.

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

# Update myspec files in current directory
myspec update

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

The registry of installed projects is stored at `~/.config/myspec/registry.json` (all platforms).

## Relationship to Other Tools

- **OpenSpec**: myspec uses OpenSpec for spec-driven change management and provides a custom schema. It does not modify any OpenSpec built-in schemas or skills.
- **Superpowers**: myspec-br borrows brainstorming methodology (Socratic dialogue, HARD-GATE, section-by-section approval) from the Superpowers project. No code or API dependency.

## Adopting myspec in Existing Projects

### New project (no OpenSpec)

```bash
myspec install /path/to/project
```

Creates OpenSpec structure, installs skills and schema. Start with `/myspec-br`.

### Existing OpenSpec project (using `spec-driven`)

```bash
myspec install /path/to/project
```

**What changes:** Default schema switches to `myspec-driven` for new changes. Existing changes are unaffected (each change records its own schema).

**What you get:** `myspec-br` for structured brainstorming, `myspec-gwt` for worktree creation, `myspec-apply` for task implementation with auto-commits, `myspec-verify` for user acceptance, `myspec-merge` for merge orchestration.

**Risk:** Low. Old changes continue using `spec-driven`. New changes use `myspec-driven`.

### Existing project using Superpowers

**Caution.** myspec-br and `superpowers:brainstorming` overlap. Both appear in `.claude/skills/`.

**Options:**
1. **Use both side by side** — manually choose `/myspec-br` or `/superpowers:brainstorming` per feature
2. **Replace superpowers brainstorming** — uninstall superpowers plugin, use myspec-br instead
3. **Keep superpowers, skip myspec** — if you use the full superpowers pipeline (brainstorming → writing-plans → subagent-driven-development → finishing), myspec adds no value

**Do not switch schema** if using `superpowers-bridge`. The DAG structures differ (superpowers-bridge includes `plan`, `retrospective`; myspec-driven includes `verify`).

### Summary

| Scenario | Install skills? | Switch schema? |
|---|---|---|
| New project | Yes | Yes (default) |
| Existing spec-driven | Yes | Yes (safe) |
| Existing superpowers-bridge | Optional | No |
| Using full superpowers pipeline | No | No |

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
- myspec-apply skill (task implementation with auto-commits)
- myspec-verify skill (user acceptance + iteration)
- myspec-merge skill (main sync + merge method + archive cleanup)
- myspec-driven custom schema with 6 artifact DAG
- Schema template files (6 templates)
- Go CLI (install / update / list / uninstall / check / doctor)
- Dynamic skill embedding (new skills auto-included in builds)
- Multi-platform binary releases (macOS, Linux, Windows)

Planned:
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
