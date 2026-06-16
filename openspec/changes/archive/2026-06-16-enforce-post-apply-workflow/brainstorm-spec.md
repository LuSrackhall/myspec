# enforce-post-apply-workflow

## Context

在 myspec 工作流的实际使用中，发现 agent 在工作树中执行 `/opsx:apply` 完成后，存在以下偏离正确流程的行为：

1. tasks.md 中的 checkbox 未被标记为完成（`- [ ]` → `- [x]`）
2. 跳过 `/opsx:verify`，直接跑 build 构建
3. 构建通过后未经用户确认，直接 `git merge` 到主分支
4. 合并后未清理工作树，而是回到工作树中给 tasks 打对钩
5. 再次合并到主分支
6. 仍然未删除工作树

**正确流程应该是：**
```
apply → verify → 用户验收 → merge（需用户确认）→ archive（在 main 上）→ worktree remove
```

**根因分析：**

1. `schema.yaml` 的 `apply.instruction` 只有 3 行，过于简略，没有 HARD-GATE
2. `myspec-br` 第 9 步的"下一步指引"是建议性命令列表，不是强制执行序列
3. 工作流中缺少用户验收暂停点
4. 没有定义 verify 失败/用户验收不通过后的迭代路径
5. 如果用户直接跑 `/opsx:apply`（不经过 myspec-br），则完全没有任何流程约束

## Goals / Non-Goals

**Goals:**

- 在 schema 层面建立 HARD-GATE，防止 agent 跳过关键步骤
- 定义完整的用户验收暂停点，确保合并前问题已被解决
- 定义用户验收不通过时的迭代路径（5 种策略）
- 确保无论用户从哪个入口进入工作流，流程约束都生效
- 在 myspec-br 第 9 步建立收尾阶段的强制执行序列

**Non-Goals:**

- 不修改 OpenSpec 通用技能文件（openspec-apply-change、openspec-verify-change、openspec-archive-change）
- 不添加新的技能或命令
- 不修改 OpenSpec CLI 代码
- 不改变现有 artifact DAG 结构

## Decisions

### Decision 1: 三层防护架构

采用三层防护确保工作流不被跳过：

| 层 | 位置 | 作用域 | 为什么有效 |
|---|------|--------|-----------|
| 第 1 层 | `schema.yaml` → `apply.instruction` | 实施阶段 | agent 在实施时直接读到，不可能绕过。即使用户直接跑 `/opsx:apply` 也会生效 |
| 第 2 层 | `templates/tasks.md` | 任务模板 | 每次查看 tasks 文件都会看到工作流提示 |
| 第 3 层 | `myspec-br` 第 9 步 | 收尾阶段 | 覆盖 merge → archive → cleanup 的强制执行序列 |

**为什么不用 hooks（settings.json）：** hooks 只能执行 shell 命令，无法表达复杂的条件逻辑（如"verify 是否完成"）。prompt 层面的 HARD-GATE 对 LLM agent 更有效。

### Decision 2: 工作树闭环原则

**工作树是一个闭环：只出不进，成品才能合并到 main。**

```
apply → verify → build/test → 用户验收
         ↑_________________________↓ 不通过：迭代修复
                                    通过
                                     ↓
                              反哺文档（确保 artifacts 匹配最终代码）
                                     ↓
                              用户确认合并
                                     ↓
                            merge → archive → worktree remove
```

合并到 main 之前，所有问题必须在工作树中解决。不接受"先合并 80%，剩余 20% 开新 change"的方案。

### Decision 3: 用户验收暂停点

在 verify 通过之后、merge 之前，agent 必须暂停并等待用户明确确认。

**验收流程：**
1. agent 展示变更摘要（修改了什么、verify 结果、build 状态）
2. agent 询问用户："这个变更是否满足你的需求？"
3. 只有用户明确确认后，才允许继续到 merge 步骤
4. 用户不通过 → 进入迭代循环（Decision 4）

### Decision 4: 迭代策略（5 种选项）

当用户验收不通过时，agent 通过分析后推荐策略，用户拥有最终选择权。

| # | 策略 | 适用场景 | 代码处理 | 文档处理 |
|---|------|---------|---------|---------|
| 1 | **原地修复**（默认推荐） | 实现细节有误、遗漏边界条件 | 保留，直接修改 | 最后反哺 |
| 2 | **同工作树开新 change** | 需要重新规划，现有代码可参考 | 保留（旧 change 代码仍在） | 新 change 独立 artifacts |
| 3 | **git reset + stash 参考** | 需要干净 baseline 但想保留参考 | stash 保存旧代码 | 从干净状态重写 |
| 4 | **git reset（完全重来）** | 彻底推翻原有方案 | 丢弃 | 从干净状态重写 |
| 5 | **废弃变更** | 放弃当前变更 | 丢弃 | 无 |

**策略 2 的注意事项：**
- 同一工作树中会存在两个 change 目录，需用不同名称
- 归档时使用 `/opsx:bulk-archive` 处理冲突
- 新 proposal 的 baseline 是当前 worktree 状态（包含第一次实施的代码），不是 main

**策略 3 的注意事项：**
- 旧代码通过 `git stash` 保留，新实施时可参考
- 需要手动清理旧 change 目录

**策略 4 的注意事项：**
- `git reset --hard` 回到实施前的 commit
- 旧代码在 git 历史中仍可恢复

**策略 5 的操作：**
```bash
git worktree remove .worktrees/change/<name>
git branch -D change/<name>
```

**Agent 职责：** 分析问题根因 → 推荐策略（附理由）→ 用户确认或选择其他。

### Decision 5: schema.yaml apply.instruction 增强

将当前 3 行的 `apply.instruction` 扩展为完整的 HARD-GATE 指令：

```yaml
apply:
  requires:
    - tasks
  tracks: tasks.md
  instruction: |
    Read context files (brainstorm-spec.md, design.md, specs/), work
    through pending tasks in tasks.md, mark each complete (`- [x]`) AS
    YOU FINISH IT — do NOT defer checkbox marking to the end.

    Pause if you hit blockers or need clarification.

    After all tasks are done:

    HARD-GATE: You MUST run /opsx:verify next. Do NOT run build, test,
    git merge, git commit to main, or any other action before verify
    completes.

    After verify completes:

    HARD-GATE: You MUST present a change summary to the user and ask
    for explicit acceptance. Do NOT proceed to merge without user
    confirmation. The summary MUST include:
    - What was changed (files, code)
    - Verify result (PASS / PASS WITH WARNINGS / FAIL)
    - Build/test status (if applicable)

    If the user does NOT accept:

    HARD-GATE: Do NOT merge. You MUST enter the iteration loop:
    1. Analyze what went wrong
    2. Recommend an iteration strategy:
       - Fix in place (default for minor issues)
       - New change in same worktree (for re-planning)
       - Git reset + stash reference (for clean baseline with reference)
       - Git reset (full redo)
       - Abandon change
    3. Present recommendation with reasoning
    4. User chooses (may pick a different strategy)
    5. Execute chosen strategy
    6. Return to apply → verify → user acceptance loop

    If the user ACCEPTS:

    HARD-GATE: Before merge, ensure ALL artifacts match the final
    implementation. Update brainstorm-spec.md, proposal.md, specs/,
    design.md, and tasks.md to reflect what was actually built
    (not just the original plan).

    Then instruct the user to:
    1. cd to repo root and git checkout main
    2. git merge change/<name> (must ask user to confirm)
    3. Run /opsx:archive on main
    4. git worktree remove .worktrees/change/<name>

    HARD-GATE: Do NOT allow any actions between merge and archive.
    After merge, the ONLY permitted next step is /opsx:archive.
    After archive, the ONLY permitted next step is git worktree remove.
```

### Decision 6: tasks.md 模板增强

在 tasks.md 模板末尾添加工作流提示：

```markdown
---

## Post-Implementation Workflow

<!-- DO NOT MODIFY THIS SECTION — it defines the required workflow after all tasks are complete -->

After completing ALL tasks above, follow this sequence strictly:

1. **Verify**: Run `/opsx:verify` to produce verify.md
2. **User Acceptance**: Present change summary, ask user to confirm the problem is solved
3. **Merge**: After user accepts, go to main branch and merge (must ask user)
4. **Archive**: Run `/opsx:archive` on main
5. **Cleanup**: `git worktree remove .worktrees/change/<name>`

**Iteration**: If user does not accept, analyze the issue and recommend:
fix in place / new change / git reset + stash / git reset / abandon.
```

### Decision 7: myspec-br 第 9 步增强

将第 9 步从"告知用户命令"改为"强制执行序列 + HARD-GATE"。

当前第 9 步输出的收尾步骤将包含以下结构（以 worktree 路径为例）：

```
HARD-GATE: 以下步骤必须严格按顺序执行，不可跳步。

1. /opsx:verify — 在工作树中完成验证
2. 展示变更摘要，询问用户是否接受
3. 用户接受后，反哺文档（确保 artifacts 匹配最终代码）
4. 回到 main: cd <repo-root> && git checkout main
5. git merge change/<name>（必须先询问用户确认，合并冲突时帮助解决或 git merge --abort）
6. /opsx:archive（在 main 上执行）

HARD-GATE: merge 完成后，禁止 build/test/编辑 task 等任何操作。
唯一的下一步是 /opsx:archive。

7. git worktree remove .worktrees/change/<name>

HARD-GATE: archive 完成后，唯一的下一步是 git worktree remove。

如果用户验收不通过，进入迭代循环：
- 分析问题 → 推荐策略（原地修/新开 change/reset+stash/reset/废弃）→ 用户选择 → 循环
```

### Decision 8: DESIGN.md 更新

在 Flow D 中加入 HARD-GATE 和迭代循环的描述。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| HARD-GATE 是 prompt 层面的，不能 100% 保证 agent 遵守 | 三层防护叠加提高概率；最高风险层（apply.instruction）在实施时直接读到 |
| 策略 2（同工作树开新 change）可能导致两个 change 目录造成选择混淆 | agent 在推荐时标注此风险；使用 `/opsx:bulk-archive` 处理 |
| 策略 3/4（git reset）涉及破坏性 git 操作 | 在 HARD-GATE 中标注这些操作需要用户确认 |
| tasks.md 模板中的 Post-Implementation Workflow 可能在归档时被误认为任务 | 用 HTML 注释标明 "DO NOT MODIFY" |
| 5 个迭代选项可能造成决策疲劳 | agent 默认推荐策略 1（原地修复），其他选项按需展示 |
