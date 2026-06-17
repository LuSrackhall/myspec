## 1. myspec-apply 技能

- [x] 1.1 创建 `.claude/skills/myspec-apply/SKILL.md`：包装 openspec-apply-change，调用 `openspec instructions apply --json` 获取 contextFiles 和 tasks，按 task group 实施并自动提交（约定式提交，用户语言/英文），完成后提示运行 myspec-verify

## 2. myspec-verify 技能

- [x] 2.1 创建 `.claude/skills/myspec-verify/SKILL.md`：包装 openspec-verify-change，3 阶段流程（文档验证 → 用户验收 → 迭代决策），验收通过后反哺所有 artifacts，提示运行 myspec-merge

## 3. myspec-merge 技能

- [x] 3.1 创建 `.claude/skills/myspec-merge/SKILL.md`：main 同步检查（比较本地/origin，引导用户决策），在工作树中 merge main 解决冲突
- [x] 3.2 myspec-merge 合并方式选择：AskUserQuestion 展示 3 种合并方式（英文写死 + 翻译注释），按用户选择执行
- [x] 3.3 myspec-merge 收尾：合并后调用 openspec-archive-change，清理工作树和分支

## 4. Schema 和文档更新

- [x] 4.1 精简 `openspec/schemas/myspec-driven/schema.yaml` 的 apply.instruction：从 ~70 行减到 ~10 行，只保留任务实施指引 + 提示运行 myspec-verify
- [x] 4.2 同步 embed 目录：将 myspec-apply、myspec-verify、myspec-merge 的 SKILL.md 复制到 `embed/skills/` 和 `internal/emb/embed/skills/`
- [x] 4.3 更新 `DESIGN.md`：Flow D 反映新的 apply → verify → merge 流程
- [x] 4.4 更新 `CLAUDE.md`：工作流技能表新增 myspec-apply、myspec-verify、myspec-merge

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
