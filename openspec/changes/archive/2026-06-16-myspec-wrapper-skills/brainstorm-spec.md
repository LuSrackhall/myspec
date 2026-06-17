# myspec-wrapper-skills

## Context

myspec 当前直接使用 OpenSpec 的通用技能（openspec-apply-change、openspec-verify-change、openspec-archive-change）来管理实施-验证-归档流程。这些技能有以下局限：

1. **apply 技能不管理 git 提交** — agent 在工作树中实施代码但从不提交，导致工作区混乱
2. **verify 技能是只读审计** — 只输出报告，不处理"用户不满意"的迭代场景
3. **没有 main 分支同步机制** — 合并前不检查 main 是否最新，可能产生冲突
4. **没有合并方式选择** — 只有 merge commit，无法选择 squash 或 rebase
5. **schema.yaml 的 apply.instruction 混合了所有流程** — 70 行指令 agent 只读一次，后续流程被遗忘

**核心矛盾：** 把所有防护放在 apply.instruction 一个地方，但 agent 不会在整个生命周期中反复读它。

## Goals / Non-Goals

**Goals:**
- 创建 3 个 myspec 包装技能（apply、verify、merge），替代 OpenSpec 命令处理核心流程
- apply 技能：每 task group 自动提交（约定式提交，用户语言/英文）
- verify 技能：文档验证 + 用户验收 + 迭代决策 + 文档反哺
- merge 技能：main 同步检查 + 合并方式选择 + 合并 + 归档 + 清理
- 精简 schema.yaml 的 apply.instruction
- 同步 embed 目录，确保 `myspec install` 分发新技能

**Non-Goals:**
- 不修改 OpenSpec 的任何技能文件（openspec-apply-change、openspec-verify-change、openspec-archive-change）
- 不创建 `/myspec:*` 斜杠命令
- 不修改 OpenSpec CLI 代码
- 不包装 propose/ff/new 等非核心 OpenSpec 技能

## Decisions

### Decision 1：包装而非替换

myspec 包装技能内部调用 OpenSpec CLI（`openspec instructions apply`、`openspec status`）获取数据，用 myspec 自己的流程逻辑包装。OpenSpec 技能文件保持原样。

**为什么不用修改 OpenSpec 技能的方式：** OpenSpec 是通用工具，不应为 myspec 的特定需求而修改。包装方式保持了关注点分离。

### Decision 2：3 个技能覆盖核心流程

| 技能 | 替代的 OpenSpec 命令 | 新增能力 |
|------|---------------------|---------|
| myspec-apply | /opsx:apply | task group 自动提交 |
| myspec-verify | /opsx:verify | 用户验收 + 迭代决策 + 文档反哺 |
| myspec-merge | （新建） | main 同步 + 合并方式 + 归档 + 清理 |

propose/ff/new/explore/continue/sync 等继续使用 OpenSpec 命令。

### Decision 3：Git 提交策略

**工作树分支：** AI 自动提交，每完成一个 task group（tasks.md 中的 `## N.` 分组）提交一次。

**提交格式：**
- 必须使用约定式提交（`feat(scope): message`、`fix(scope): message` 等）
- commit message 使用用户偏好语言，无偏好时默认英文

**main 分支：** 所有操作（pull/push/merge）必须用户决策。

### Decision 4：合并方式选择

提供 3 种合并方式，用 AskUserQuestion 展示（英文写死 + 翻译注释）：
- Create a merge commit（保留分支历史）
- Squash and merge（压缩为一个 commit）
- Rebase（线性历史）

### Decision 5：main 同步策略

合并前检查本地 main vs origin/main 的领先状态：
- origin/main 领先 → 提示用户 pull
- 本地 main 领先 → 提示用户 push
- 同步 → 跳过

用户确认后，在工作树中 `git merge main`，将冲突解决隔离在工作树中。

### Decision 6：verify 包含用户验收和迭代

verify 技能不再是只读审计，而是包含 3 个阶段：
1. 文档验证（复用 OpenSpec verify 逻辑）
2. 用户验收（展示摘要，询问用户）
3. 迭代决策（不接受时推荐 5 种策略）

验收通过后反哺所有 artifacts 匹配最终代码。

### Decision 7：schema.yaml 精简

apply.instruction 从 ~70 行精简到 ~10 行，只保留：
- 任务实施指引（读上下文、标记 checkbox）
- 完成后提示运行 myspec-verify

所有 HARD-GATE 和后续流程移至 myspec 包装技能中。

### Decision 8：通过 Skill 工具调用，无斜杠命令

myspec 包装技能通过 Skill 工具直接调用（`.claude/skills/myspec-apply/SKILL.md` 等），不创建 `/myspec:*` 斜杠命令。技能之间通过名称互相提示。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| agent 可能仍使用 /opsx:apply 而非 myspec-apply | 更新 CLAUDE.md 指引；schema.apply.instruction 提示用 myspec-verify |
| 包装技能比原始 OpenSpec 技能更长，agent 可能跳过部分指令 | 每个技能职责单一，比 70 行混合指令更容易遵守 |
| main 同步检查增加了合并前的步骤 | 这是必要的安全措施，防止合并冲突 |
| 3 种合并方式增加了用户决策负担 | 只在合并时展示一次，不影响日常开发 |
