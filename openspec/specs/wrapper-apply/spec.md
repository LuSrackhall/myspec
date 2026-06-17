## ADDED Requirements

### Requirement: apply-task-group-commits

The myspec-apply skill SHALL commit code after completing each task group, using conventional commit format with the user's preferred language.

#### Scenario: task group completed

WHEN the agent finishes all tasks in a task group (## N. section)
THEN the agent MUST run `git add -A && git commit` with a conventional commit message (`<type>(<scope>): <description>`), using the user's preferred language (default English)

#### Scenario: all tasks completed

WHEN all task groups are completed
THEN the agent MUST prompt the user to run the myspec-verify skill, and MUST NOT perform any other post-implementation action (no build, no merge, no archive)

### Requirement: apply-openspec-integration

The myspec-apply skill SHALL use OpenSpec CLI commands to retrieve task lists and context files, without modifying OpenSpec skill files.

#### Scenario: apply retrieves context

WHEN myspec-apply is invoked
THEN the skill MUST call `openspec instructions apply --change "<name>" --json` to get contextFiles and task list, then read all context files before implementation
