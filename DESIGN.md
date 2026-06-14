# myspec 设计文档

## 项目概述

myspec 是一个 Claude Code 工作流管理工具，用于：
1. 定义可复用的开发工作流技能
2. 安装/更新技能到多个项目
3. 通过 Go 二进制 + 内嵌文件实现跨项目分发

## 设计决策记录

### 决策 1：工作流流程

**最终确定的流程：**

```
1. main: 用户提需求
2. main: 用 brainstorming 方法论讨论（调用 /myspec-br）
   - 逐一提问
   - 2-3 方案对比
   - 分段展示设计
   - 用户批准设计
3. main: AskUserQuestion "创建工作树？"
4. main: EnterWorktree
5. worktree: openspec new change "<name>"
6. worktree: 写入 brainstorm-spec.md 到 change 目录
7. worktree: 调用 opsx:propose（回答"继续已有 change"）
8. worktree: opsx:apply → cargo run → opsx:verify
9. main: 询问用户确认 → git merge
10. main: opsx:archive
```

**关键约束：**
- 归档只在 main 分支进行
- 合并前必须询问用户，不允许私自合并
- 使用 git merge（非 squash），保留分支历史

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

### 决策 5：opsx:propose 与 brainstorming 的关系

- `opsx:propose` = `opsx:new` + 循环 `opsx:continue`（批量生成 artifact）
- propose 是"你说做什么我就生成什么"，没有设计讨论过程
- brainstorming 的产出（brainstorm-spec.md）作为 propose 的输入上下文，提高 artifact 质量

### 决策 6：设计文档必须提交

worktree 从 main 的最新 commit 创建。如果 propose 后不提交，worktree 里看不到设计文档。

**结论：** 设计文档必须提交（`openspec new change` 创建骨架，写入 brainstorm-spec.md，然后 propose 读取它）。

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

## 待实现

- [ ] myspec Go 工具开发（install / update / list / uninstall / check）
- [ ] SKILL.md 验证（在 bevy-test 项目中实际使用 myspec-br）
- [ ] 全局 skill 安装支持
- [ ] 项目级覆盖机制
