## Why

Agent 在工作树中执行 `/opsx:apply` 后跳过了 verify、archive、worktree cleanup 等关键步骤，直接合并到 main。根因是 `apply.instruction` 过于简略且无 HARD-GATE，`myspec-br` 第 9 步仅为建议性命令列表而非强制执行序列。工作流中缺少用户验收暂停点和迭代路径定义。

## What Changes

1. **增强 `schema.yaml` 的 `apply.instruction`**（第 1 层防护）：在实施阶段加入 HARD-GATE，要求 apply 完成后必须 verify → 用户验收 → merge → archive → cleanup，不可跳步
2. **增强 `templates/tasks.md`**（第 2 层）：在任务模板末尾添加 Post-Implementation Workflow 提示，确保查看 tasks 时总能看到后续流程
3. **增强 `myspec-br` 第 9 步**（第 3 层）：从"告知用户命令"改为"agent 强制执行序列 + HARD-GATE"
4. **定义用户验收暂停点**：verify 通过后、merge 前，agent 必须暂停等待用户确认
5. **定义 5 种迭代策略**：用户验收不通过时，agent 分析后推荐策略，用户可选：原地修复、同工作树开新 change、git reset+stash、git reset、废弃
6. **更新 `DESIGN.md`**：Flow D 加入 HARD-GATE 和迭代循环描述

## Capabilities

### New Capabilities
- `workflow-enforcement`: 三层 HARD-GATE 防护架构，用户验收暂停点，5 种迭代策略，工作树闭环原则

### Modified Capabilities

## Impact

- `openspec/schemas/myspec-driven/schema.yaml`: `apply.instruction` 字段扩展
- `openspec/schemas/myspec-driven/templates/tasks.md`: 模板末尾添加工作流提示
- `.claude/skills/myspec-br/SKILL.md`: 第 9 步重写
- `DESIGN.md`: Flow D 更新
- 不影响现有 CLI 代码、specs、或其他技能文件
