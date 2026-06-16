## 1. Schema 层防护（第 1 层）

- [x] 1.1 增强 `openspec/schemas/myspec-driven/schema.yaml` 的 `apply.instruction` 字段：添加 verify 前 HARD-GATE、用户验收 HARD-GATE、迭代循环指令、merge 后 HARD-GATE、archive 后 HARD-GATE

## 2. 模板层防护（第 2 层）

- [x] 2.1 修改 `openspec/schemas/myspec-driven/templates/tasks.md`：在模板末尾添加 `## Post-Implementation Workflow` 区块，包含 5 步收尾流程、迭代策略概要、HTML 注释 `DO NOT MODIFY` 标记

## 3. 技能层防护（第 3 层）

- [x] 3.1 重写 `.claude/skills/myspec-br/SKILL.md` 第 9 步"告知用户下一步"：从建议性命令列表改为带 HARD-GATE 的强制执行序列，包含用户验收暂停点和迭代循环指引
- [x] 3.2 更新 `.claude/skills/myspec-br/SKILL.md` 第 9 步的无工作树路径：添加简化版 HARD-GATE 序列（verify → 验收 → archive）

## 4. 文档更新

- [x] 4.1 更新 `DESIGN.md` 的 Flow D：在 apply → merge 之间插入 verify + 用户验收暂停点 + HARD-GATE 标记
- [x] 4.2 更新 `DESIGN.md` 的 Flow D：添加迭代循环说明（验收不通过时的处理路径）

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
