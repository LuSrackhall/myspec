## Context

myspec 使用 myspec-driven schema 管理工作流。schema 的 `apply.instruction` 是 agent 在实施阶段读取的核心指令，但当前只有 3 行，没有 HARD-GATE。`myspec-br` 的第 9 步"告知用户下一步"是建议性命令列表，agent 经常忽略。需要在不修改 OpenSpec 通用技能文件的前提下，通过 schema 模板和 myspec 专属技能文件建立强制执行机制。

## Goals / Non-Goals

**Goals:**
- 在 `schema.yaml` 的 `apply.instruction` 中建立 HARD-GATE（第 1 层）
- 在 `templates/tasks.md` 中添加工作流提示（第 2 层）
- 在 `myspec-br` 第 9 步建立强制执行序列（第 3 层）
- 定义用户验收暂停点和 5 种迭代策略
- 更新 DESIGN.md 的 Flow D

**Non-Goals:**
- 不修改 OpenSpec 通用技能文件
- 不添加新的 CLI 命令或技能
- 不修改 OpenSpec CLI 代码

## Decisions

### Decision 1: apply.instruction 扩展策略

将当前 3 行的 `apply.instruction` 扩展为包含完整 HARD-GATE 的指令块。

**改动位置：** `openspec/schemas/myspec-driven/schema.yaml` 的 `apply.instruction` 字段

**新增内容：**
- `HARD-GATE: verify 前禁止 build/merge` — 防止实施后直接跑 build 或 merge
- `HARD-GATE: 用户验收前禁止 merge` — 防止跳过用户确认
- 迭代循环指令 — 5 种策略的描述和选择流程
- `HARD-GATE: merge 后禁止除 archive 外的任何操作` — 防止 merge 后做 build/re-merge
- `HARD-GATE: archive 后唯一的下一步是 worktree remove` — 防止遗留工作树

**为什么放在 apply.instruction 而不是其他位置：**
- apply.instruction 在 agent 实施任务时直接读取，时机最恰当
- 即使用户直接跑 `/opsx:apply`（跳过 myspec-br），也会读到这些约束
- 不需要创建新技能或修改 OpenSpec 通用技能

### Decision 2: tasks.md 模板添加工作流提示

在 `templates/tasks.md` 模板末尾添加 `## Post-Implementation Workflow` 区块，用 HTML 注释标记 `DO NOT MODIFY`。

**内容包含：**
- 5 步收尾流程（verify → 验收 → merge → archive → cleanup）
- 迭代提示（验收不通过时的 5 种策略概要）

**设计考量：**
- 此区块不是任务，不会被 apply 解析为 checkbox
- 每次查看 tasks 文件都会看到，作为持续提醒
- HTML 注释标记防止被误修改或删除

### Decision 3: myspec-br 第 9 步重写

将当前的"告知用户命令"改为"agent 强制执行序列"。

**当前行为：** 输出 shell 命令片段，建议用户手动执行
**目标行为：** 输出带 HARD-GATE 的强制序列，agent 引导用户逐步执行

**关键变化：**
- 添加 `HARD-GATE` 标记到每个关键步骤之间
- 从"建议性"改为"强制性"语气
- 补充用户验收暂停点指引
- 补充迭代循环指引

### Decision 4: DESIGN.md Flow D 更新

在 Flow D 的第 8-11 步之间插入：
- verify 步骤（已有但未显式强调）
- 用户验收暂停点
- HARD-GATE 标记
- 迭代循环说明

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| HARD-GATE 是 prompt 层面的，LLM 不保证 100% 遵守 | 三层叠加提高概率；第 1 层（apply.instruction）时机最关键 |
| tasks.md 模板中 Post-Implementation Workflow 可能被归档工具误处理 | HTML 注释标记 + 明确的 section 标题 |
| 5 个迭代选项可能造成决策疲劳 | agent 默认推荐策略 1，其余按需展示 |
| apply.instruction 扩展后可能影响其他使用 myspec-driven schema 的项目 | 不修改通用技能，只改 schema 模板；扩展是指令性文本，不影响解析 |
