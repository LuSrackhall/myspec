# myspec 项目宪法

## 项目概述

myspec 是一个 Claude Code 工作流管理工具，用于定义可复用的开发工作流技能并通过 Go CLI 分发到多个项目。

## 自举策略（Alternative E）

**核心原则：在 myspec 仓库自身中完成工作流验证和 Go CLI 开发。**

```
Phase 1: 验证工作流 + 构建 CLI（同一个开发周期）
  1. openspec init --tools claude（在本仓库中）
  2. 设置 schema: myspec-driven
  3. /myspec-br 脑暴 Go CLI 设计
  4. /opsx:propose 生成 artifacts
  5. /opsx:apply 实现 Go CLI
  6. /opsx:verify 验证
  7. 归档

Phase 2: 用 CLI 安装到外部项目测试
  8. go build -o myspec .
  9. myspec install /path/to/test-project
  10. 验证安装结果

Phase 3: 迭代
  11. 修复问题，用 myspec update 推送
```

**为什么自举：**
- Go CLI 是分发机制，不影响工作流正确性
- 工作流可以在本仓库中独立验证
- Dogfooding：用自己的工具构建自己的工具
- 不需要测试项目，不需要复制回文件，不丢失 git 历史

## 工作流技能

| 技能 | 位置 | 职责 |
|---|---|---|
| myspec-br | `.claude/skills/myspec-br/SKILL.md` | 脑暴编排器（Flow D） |
| myspec-gwt | `.claude/skills/myspec-gwt/SKILL.md` | git worktree 创建 |
| myspec-apply | `.claude/skills/myspec-apply/SKILL.md` | 任务实施 + task group 自动提交 |
| myspec-verify | `.claude/skills/myspec-verify/SKILL.md` | 文档验证 + 用户验收 + 迭代决策 |
| myspec-merge | `.claude/skills/myspec-merge/SKILL.md` | main 同步 + 合并方式选择 + 归档清理 |

## 自定义 Schema

| 文件 | 说明 |
|---|---|
| `openspec/schemas/myspec-driven/schema.yaml` | 6 artifact DAG |
| `openspec/schemas/myspec-driven/templates/` | 6 个模板文件 |

**DAG：** brainstorm-spec → proposal → specs + design → tasks → [apply] → verify

## 设计文档

`DESIGN.md` 包含所有设计决策。修改工作流前必须先更新 DESIGN.md。

## Git 提交规范

所有 commit message **必须使用中文**。格式：

```
feat(simulation): 新增朝向系统组件
fix(render): 修复盾牌血条显示问题
docs: 更新设计文档
```

## 关键约束

- 归档在 main 上进行（基于权威 specs），合并后立即归档并提交
- 合并前必须询问用户，不允许私自合并
- 使用 git merge（非 squash），保留分支历史
- myspec-br 不自动调用 propose，只告知用户下一步
- myspec-br 不判断 OpenSpec 是否可用（CLI 负责安装/检测）
