## ADDED Requirements

### Requirement: verify-user-acceptance

The myspec-verify skill SHALL present a change summary and wait for explicit user acceptance before proceeding.

#### Scenario: verification completes

WHEN document verification (Completeness/Correctness/Coherence) completes
THEN the agent MUST present a change summary including: files changed, verify result (PASS/WARNINGS/FAIL), and key changes, then ask the user for explicit acceptance

#### Scenario: user accepts

WHEN the user explicitly accepts the change
THEN the agent MUST backfill all artifacts (brainstorm-spec, proposal, specs, design, tasks) to match the final implementation, then prompt the user to run the myspec-merge skill

#### Scenario: user does not accept

WHEN the user does NOT accept the change
THEN the agent MUST analyze the root cause, recommend an iteration strategy with reasoning (fix in place / new change / git reset + stash / git reset / abandon), present the recommendation, let the user choose, execute the chosen strategy, and return to myspec-apply

### Requirement: verify-openspec-integration

The myspec-verify skill SHALL use OpenSpec CLI commands to retrieve context files and perform document verification, without modifying OpenSpec skill files.

#### Scenario: verify retrieves context

WHEN myspec-verify is invoked
THEN the skill MUST call `openspec instructions apply --change "<name>" --json` to get contextFiles, then read all artifacts for verification
