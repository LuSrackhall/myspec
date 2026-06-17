# myspec 设计文档

## 项目概述

myspec 是一个 Claude Code 工作流管理工具，用于：
1. 定义可复用的开发工作流技能
2. 安装/更新技能到多个项目
3. 通过 Go 二进制 + 内嵌文件实现跨项目分发

## 设计决策记录

### 决策 1：工作流流程

**最终确定的流程（Flow D）：**

```
1. main: 用户提需求
2. main: myspec-br 脑暴（逐一提问、方案对比、分段设计、自检、用户审阅）
3. main: myspec-br 询问用户"是否创建工作树？"
   ├─ 是 → myspec-br 调用 myspec-gwt 创建 worktree
   │        myspec-br cd 到 worktree
   └─ 否 → myspec-br 留在 main（警告会污染 main）
4. myspec-br: openspec new change "<name>"
5. myspec-br: 写入 brainstorm-spec.md 到 change 目录，提交
6. myspec-br: 告知用户"运行 /opsx:propose 继续"
7. 用户: 运行 /opsx:propose（或 /opsx:ff）
8. worktree: myspec-apply（实施任务，每 task group 自动提交）
9. worktree: myspec-verify（文档验证 → 用户验收 → 迭代决策）
   ├─ 用户接受 → 反哺文档 → 继续步骤 10
   └─ 用户不接受 → 迭代（分析 → 推荐策略 → 用户选择 → 回到步骤 8）
10. main: myspec-merge
    a. main 同步检查（比较本地/origin main，引导用户决策）
    b. 在工作树中 merge main（冲突在工作树中解决）
    c. 选择合并方式（merge commit / squash / rebase）
    d. 执行合并（用户确认）
    e. opsx:archive → git commit
    f. git worktree remove + git branch -d
```
```

**myspec-br 是编排器：** 它拥有脑暴上下文，决定何时创建工作树、何时创建 change 目录、何时写入文件。myspec-gwt 和 propose 是它委托的工具。

**收尾流程设计：**
- 合并后归档（merge → archive → commit → worktree remove），不是归档后合并
- 原因：archive 的 spec sync 必须基于 main 的权威 `openspec/specs/`，worktree 的 specs 可能过时
- archive 产生未提交的 spec 编辑和目录移动，必须在 archive 后立即 commit
- worktree remove 是最后一步，确保 archive 成功后再清理

**关键约束：**
- 归档在 main 上进行（基于权威 specs），合并后立即归档并提交
- 合并前必须询问用户，不允许私自合并
- 使用 git merge（非 squash），保留分支历史
- myspec-br 不自动调用 propose，只告知用户下一步
- apply 完成后必须 verify → 用户验收，不可跳步（HARD-GATE）
- merge 后唯一的下一步是 archive（HARD-GATE）
- archive 后唯一的下一步是 worktree remove（HARD-GATE）
- 用户验收不通过时进入迭代循环（5 种策略），不直接合并

### 决策 2：为什么 brainstorming 在 main 上，propose 在 worktree 上

**讨论阶段（main）：**
- 只产生对话，不产生文件
- 没有并发风险（暂存区无冲突）
- 可以查看完整代码上下文

**设计文档 + 实施（worktree）：**
- 避免废弃提案污染 main
- worktree 废弃时带走一切
- 隔离不同需求的实施

**worktree 创建时机：讨论结束后、propose 前。** 因为：
- 讨论可能否决需求，此时创建 worktree 是浪费
- 讨论阶段不产生文件，没有并发风险
- 所有文件操作都在 worktree 里进行，main 保持干净

### 决策 3：为什么不直接调用 superpowers:brainstorming

brainstorming 的终止状态硬编码为 `writing-plans`（superpowers 的计划生成技能）。调用 brainstorming 会自动链入 superpowers 的完整流程（writing-plans → subagent-driven-development → finishing-a-development-branch），无法配置关闭。

**解决方案：** 从 GitHub 仓库 (https://github.com/obra/superpowers) 拿 brainstorming 的方法论，改写为 OpenSpec artifact 驱动的版本。保留核心方法论（HARD-GATE、逐一提问、方案对比、分段批准、自校审），修改产出路径和终止状态。

### 决策 4：为什么 opsx:explore 不够用

| | opsx:explore | brainstorming |
|---|---|---|
| 定位 | 调查现状 | 从模糊想法产出完整设计 |
| 产出 | 无（纯对话） | 设计文档 |
| 是否系统性提问 | 偶尔澄清 | 必须逐一提问 |
| 是否给方案 | 不给 | 2-3 方案 + 推荐 |
| 是否分段批准 | 不需要 | 必须 |

两者不重叠。explore 是"看看发生了什么"，brainstorming 是"一起想清楚要做什么"。

### 决策 5：自定义 OpenSpec Schema `myspec-driven`

- 创建自定义 schema（fork 自 `spec-driven`），将 `brainstorm-spec` 注册为一等 artifact
- `brainstorm-spec` → `proposal` → `specs` + `design` → `tasks` → `[apply]` → `verify`
- `proposal` 的 instruction 从 brainstorm-spec.md 提取内容，而非从零生成
- myspec-br 是独立技能，可脱离 OpenSpec 使用；当 OpenSpec 可用时自动集成到 schema 流程
- 不修改 OpenSpec 任何内置 schema 或技能

### 决策 6：设计文档在 worktree 中提交

设计文档（brainstorm-spec.md）在 worktree 中创建和提交，不在 main 上。

**流程：** myspec-br 脑暴完成后 → 创建 worktree → cd 到 worktree → `openspec new change` → 写入 brainstorm-spec.md → 提交。

**与 Decision 2 的一致性：** main 上不产生任何文件，所有文件操作在 worktree 中进行。废弃时删除 worktree 即可，main 无痕迹。

**无 worktree 降级：** 如果用户选择不创建工作树，设计文档在 main 上创建，myspec-br 会警告用户"废弃时需手动清理"。

### 决策 7：废弃方案处理

- worktree 里的废弃方案：删除 worktree 即可，main 无痕迹
- 已合并到 main 的废弃方案：通过 `opsx:archive` 清理
- 不需要 revert 提交污染 main 历史

### 决策 8：并发安全

- 讨论阶段（main）：只对话不产生文件，无并发风险
- 实施阶段（worktree）：各自隔离，互不干扰
- 设计文档在 worktree 里：不同 worktree 有独立暂存区，无并发提交问题
- 两个 worktree 改同一文件：合并时可能出现 merge conflict，这是正常的 git 行为

### 决策 9：磁盘空间

- Rust 编译产物大（target/ 约 10GB）
- 计划将编译移到 GitHub Actions
- 本地只用 cargo run 手动测试
- 共享 CARGO_TARGET_DIR 可省空间但会导致增量缓存失效
- 当前单人开发场景下，独立 target 目录够用

### 决策 10：工具架构

**名称：** myspec

**技术栈：** Go 二进制 + 内嵌文件（`go:embed`）

**CLI 接口：**
```bash
myspec install <project-path>     # 安装技能到项目
myspec update                     # 更新所有已安装项目
myspec update <project-path>      # 更新指定项目
myspec list                       # 列出已安装项目及版本
myspec uninstall <project-path>   # 从项目卸载
myspec check                      # 检查是否有新版本可用
```

**安装方式：** 文件复制（非符号链接），项目自包含，git 可追踪

**更新策略：** 全量替换旧版本

**注册表：** `~/.config/myspec/registry.json` 记录已安装项目路径和版本

**技能命名空间：** `myspec-` 前缀

**当前技能：** `myspec-br`（需求讨论与设计）

### 决策 11：与 OpenSpec 和 Superpowers 的关系

- **不修改** OpenSpec 的任何技能/命令
- **不修改** Superpowers 的任何技能
- **借鉴** Superpowers brainstorming 的方法论（对话技巧，非 API 依赖）
- **调用** OpenSpec 的 CLI 和技能（`openspec new change`、`opsx:propose`）
- myspec 技能是**独立的入口**，在讨论阶段替代 `opsx:explore`

### 决策 12：skill 文件放置

**项目级：** `.claude/skills/myspec-br/SKILL.md`

**全局级（未来）：** `~/.claude/skills/myspec-br/SKILL.md`

**跨项目复用：** 通过 myspec Go 工具安装/更新

**项目级覆盖：** 项目 `.claude/skills/` 下的同名文件优先于全局安装的

## 已完成

- [x] myspec-br SKILL.md（脑暴技能，Flow D 编排器：脑暴 → 工作树决策 → change → 写入 → 告知用户）
- [x] myspec-gwt SKILL.md（git worktree 创建，支持 main/master、损坏检测、sandbox 降级）
- [x] myspec-driven 自定义 OpenSpec Schema（brainstorm-spec → proposal → specs + design → tasks → verify）
- [x] Schema 模板文件（brainstorm-spec.md、verify.md + 继承的 proposal/design/spec/tasks）

## 待实现

- [ ] myspec Go 工具开发（install / update / list / uninstall / check）
- [ ] 全局 skill 安装支持
- [ ] 项目级覆盖机制

**收尾流程：** 不需要独立技能。收尾指令内联在 myspec-br step 9 中（merge → archive → worktree remove）。详见 myspec-br SKILL.md。

## 与其他工具的职责边界

- **OpenSpec**：规范驱动的变更管理，myspec 通过自定义 schema 深度集成，不修改其内置 schema
- **superpowers**：借鉴 brainstorming 方法论，不依赖其代码、API 或插件
- **myspec**：开发工作流管理，为 Claude Code 扩展而生
