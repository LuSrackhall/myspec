### Requirement: apply-instruction-hard-gates

The myspec-driven schema's `apply.instruction` SHALL contain HARD-GATE markers that enforce the post-apply workflow sequence. The agent MUST NOT perform build, test, git merge, or any other action between the following steps:

#### Scenario: agent completes all tasks

WHEN the agent finishes all tasks in tasks.md
THEN the agent MUST run `/opsx:verify` next and MUST NOT run build, test, git merge, or any other action before verify completes

#### Scenario: verify completes and user acceptance required

WHEN `/opsx:verify` completes
THEN the agent MUST present a change summary to the user and ask for explicit acceptance before proceeding. The summary SHALL include: what was changed, verify result, and build/test status

#### Scenario: user accepts the change

WHEN the user explicitly accepts the change
THEN the agent MUST ensure all artifacts match the final implementation, then instruct the user to: cd to repo root, git checkout main, git merge (with user confirmation), run /opsx:archive, and git worktree remove

#### Scenario: user does not accept the change

WHEN the user does NOT accept the change
THEN the agent MUST enter the iteration loop: analyze the issue, recommend an iteration strategy (fix in place / new change / git reset + stash / git reset / abandon), present recommendation with reasoning, let the user choose, execute the chosen strategy, and return to the apply → verify → user acceptance loop

#### Scenario: merge completes

WHEN `git merge` completes successfully
THEN the agent MUST NOT perform any action except `/opsx:archive`. Build, test, task editing, re-merge, and all other operations are forbidden

#### Scenario: archive completes

WHEN `/opsx:archive` completes successfully
THEN the agent MUST NOT perform any action except `git worktree remove`. All other operations are forbidden

### Requirement: tasks-template-workflow-hint

The myspec-driven schema's `tasks.md` template SHALL contain a `## Post-Implementation Workflow` section at the end with HTML comment markers indicating `DO NOT MODIFY`.

#### Scenario: tasks template includes workflow hint

WHEN the tasks.md template is rendered
THEN the file SHALL end with a `## Post-Implementation Workflow` section containing: the 5-step post-implementation sequence (verify → user acceptance → merge → archive → cleanup), iteration strategy overview, and HTML comment markers preventing modification

### Requirement: myspec-br-enforced-post-steps

The myspec-br skill's step 9 ("Inform User of Next Steps") SHALL contain HARD-GATE markers enforcing the post-apply workflow sequence.

#### Scenario: myspec-br reaches step 9 with worktree path

WHEN myspec-br completes the design phase and the user chose worktree isolation
THEN step 9 SHALL output a mandatory sequence with HARD-GATE markers between each step: /opsx:verify → user acceptance → docs backfill → git checkout main → git merge (ask user) → /opsx:archive → git worktree remove

#### Scenario: myspec-br reaches step 9 without worktree path

WHEN myspec-br completes the design phase and the user chose not to create a worktree
THEN step 9 SHALL output a mandatory sequence with HARD-GATE markers: /opsx:verify → user acceptance → docs backfill → /opsx:archive

#### Scenario: myspec-br step 9 includes iteration guidance

WHEN myspec-br step 9 is generated
THEN it SHALL include iteration loop guidance: when user acceptance fails, the agent MUST analyze the issue and recommend one of 5 strategies (fix in place / new change in same worktree / git reset + stash / git reset / abandon), with the user having final choice
