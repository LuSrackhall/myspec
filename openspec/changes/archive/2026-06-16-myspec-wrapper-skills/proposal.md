## Why

myspec 直接使用 OpenSpec 通用技能处理核心流程，存在 5 个问题：apply 不管理 git 提交、verify 是只读审计不处理迭代、没有 main 同步机制、没有合并方式选择、schema 的 apply.instruction 混合所有流程导致 agent 遗忘后续步骤。需要创建 myspec 包装技能来补充这些能力。

## What Changes

1. **新建 myspec-apply 技能** — 包装 openspec-apply-change，增加 task group 自动提交（约定式提交，用户语言/英文）
2. **新建 myspec-verify 技能** — 包装 openspec-verify-change，增加用户验收 + 迭代决策 + 文档反哺
3. **新建 myspec-merge 技能** — 新建，处理 main 同步检查 + 合并方式选择（3 种）+ 合并 + 归档 + 清理
4. **精简 schema.yaml 的 apply.instruction** — 从 ~70 行减到 ~10 行，后续流程移至包装技能
5. **同步 embed 目录** — 新增的 3 个技能加入 embed/，通过 myspec install 分发
6. **更新 DESIGN.md 和 CLAUDE.md** — 反映新的工作流架构

## Capabilities

### New Capabilities
- `wrapper-apply`: 包装 OpenSpec apply 技能，增加 task group 约定式提交和后续流程提示
- `wrapper-verify`: 包装 OpenSpec verify 技能，增加用户验收、迭代决策和文档反哺
- `wrapper-merge`: 新建收尾技能，处理 main 同步、合并方式选择、合并执行、归档和工作树清理

### Modified Capabilities

## Impact

- 新增 `.claude/skills/myspec-apply/SKILL.md`
- 新增 `.claude/skills/myspec-verify/SKILL.md`
- 新增 `.claude/skills/myspec-merge/SKILL.md`
- 修改 `openspec/schemas/myspec-driven/schema.yaml`（apply.instruction 精简）
- 同步 `embed/skills/` 和 `internal/emb/embed/skills/`
- 更新 `DESIGN.md` Flow D
- 更新 `CLAUDE.md` 工作流技能表
- 不影响 OpenSpec 技能文件和 CLI 代码
