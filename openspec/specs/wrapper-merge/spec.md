## ADDED Requirements

### Requirement: merge-main-sync

The myspec-merge skill SHALL check if local main and origin/main are in sync, and guide the user to resolve any divergence before merging.

#### Scenario: origin/main is ahead

WHEN origin/main has commits that local main does not have
THEN the agent MUST inform the user and ask whether to pull (user decision). After user confirms, the agent MUST update local main, then merge main into the worktree branch to resolve conflicts in the worktree

#### Scenario: local main is ahead

WHEN local main has commits that origin/main does not have
THEN the agent MUST inform the user and ask whether to push (user decision)

#### Scenario: main is in sync

WHEN local main and origin/main point to the same commit
THEN the agent MUST skip the sync step and proceed to merge method selection

### Requirement: merge-method-selection

The myspec-merge skill SHALL present three merge methods and let the user choose.

#### Scenario: user selects merge method

WHEN the main sync check completes
THEN the agent MUST present three options using AskUserQuestion: Create a merge commit, Squash and merge, Rebase. Options SHALL be in English with a comment indicating the agent may translate to the user's preferred language

#### Scenario: merge commit selected

WHEN the user selects "Create a merge commit"
THEN the agent MUST execute `git merge change/<name>`

#### Scenario: squash merge selected

WHEN the user selects "Squash and merge"
THEN the agent MUST execute `git merge --squash change/<name>` followed by `git commit`

#### Scenario: rebase selected

WHEN the user selects "Rebase"
THEN the agent MUST execute `git checkout main && git rebase change/<name>`

### Requirement: merge-archive-cleanup

The myspec-merge skill SHALL archive the change and clean up the worktree after a successful merge.

#### Scenario: merge completes

WHEN the merge executes successfully
THEN the agent MUST run the openspec-archive-change skill, then remove the worktree (`git worktree remove`) and delete the branch (`git branch -d`)

### Requirement: merge-main-operations-user-decision

All main branch operations SHALL require explicit user decision. Worktree branch commits SHALL be handled automatically by the agent.

#### Scenario: main branch operation

WHEN any operation targets the main branch (pull, push, merge, rebase)
THEN the agent MUST ask the user for confirmation before executing
