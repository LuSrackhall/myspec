## Context

myspec 当前直接使用 OpenSpec 通用技能处理 apply/verify/archive 流程。这些技能缺少 git 提交管理、用户验收、迭代决策、main 同步和合并方式选择等能力。需要通过 myspec 包装技能补充这些能力，同时不修改 OpenSpec 原始文件。

## Goals / Non-Goals

**Goals:**
- 创建 myspec-apply：task group 自动提交（约定式提交，用户语言/英文）
- 创建 myspec-verify：文档验证 + 用户验收 + 迭代决策 + 文档反哺
- 创建 myspec-merge：main 同步 + 合并方式选择 + 合并 + 归档 + 清理
- 精简 schema.yaml 的 apply.instruction
- 通过 Skill 工具调用，无斜杠命令

**Non-Goals:**
- 不修改 OpenSpec 技能文件
- 不创建 `/myspec:*` 斜杠命令
- 不包装 propose/ff/new 等非核心 OpenSpec 技能

## Decisions

### Decision 1：包装架构

每个 myspec 包装技能内部调用 OpenSpec CLI 获取数据（`openspec instructions apply --json`、`openspec status --json`），用 myspec 流程逻辑包装。OpenSpec 技能文件保持原样。

### Decision 2：myspec-apply 流程

```
1. openspec instructions apply → 获取 contextFiles + tasks
2. 读取 contextFiles（brainstorm-spec、design、specs、tasks）
3. 按 task group 实施：
   a. 实施该 group 所有 tasks
   b. 标记每个 task 为 - [x]
   c. git add -A && git commit（约定式提交，用户语言/英文）
4. 完成后提示："运行 myspec-verify 技能"
```

commit message 格式：`<type>(<scope>): <description>`
- type: feat/fix/docs/refactor/test/chore
- scope: 从变更名或 task group 名推断
- description: 用户语言优先，无偏好时英文

### Decision 3：myspec-verify 流程

```
阶段 1：文档验证
- openspec instructions apply → contextFiles
- 3 维度检查：Completeness/Correctness/Coherence
- 输出验证报告

阶段 2：用户验收
- 展示变更摘要（文件变更、verify 结果、关键说明）
- AskUserQuestion："变更是否满足需求？"

阶段 3a：用户接受
- 反哺所有 artifacts（brainstorm-spec、proposal、specs、design、tasks）
- 提示："运行 myspec-merge 技能"

阶段 3b：用户不接受
- 分析根因 → 推荐策略（5 种）→ 用户选择 → 执行
- 返回 myspec-apply 重新实施
```

### Decision 4：myspec-merge 流程

```
阶段 1：main 同步检查
- git fetch origin
- 比较本地 main vs origin/main
- 引导用户决策（pull/push/跳过）
- 工作树中 git merge main（冲突在工作树中解决）

阶段 2：合并方式选择
- AskUserQuestion 展示 3 种（英文写死 + 翻译注释）：
  · Create a merge commit
  · Squash and merge
  · Rebase

阶段 3：执行
- 按用户选择执行合并
- openspec archive（调用 openspec-archive-change 技能）
- git worktree remove + git branch -d
```

### Decision 5：schema.yaml 精简

apply.instruction 从 ~70 行减到 ~10 行：
- 保留：任务实施指引 + 完成后提示 myspec-verify
- 移除：HARD-GATE、用户验收、迭代循环、merge/archive/cleanup 指引

### Decision 6：Skill 工具调用

技能通过 Skill 工具直接调用（`.claude/skills/myspec-apply/SKILL.md`），不创建斜杠命令。技能间通过名称互相提示。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| agent 仍可能用 /opsx:apply 而非 myspec-apply | schema.apply.instruction 提示用 myspec-verify；CLAUDE.md 指引 |
| 包装技能更长可能被 agent 跳过 | 每个技能职责单一，比混合指令更易遵守 |
| main 同步增加合并前步骤 | 必要的安全措施，防止合并冲突 |
