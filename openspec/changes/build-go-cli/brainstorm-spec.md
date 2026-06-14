## Context

myspec 是一个 Claude Code 工作流管理工具，用于定义可复用的开发工作流技能并通过 Go CLI 分发到多个项目。

当前已完成的组件：
- myspec-br SKILL.md（脑暴编排器）
- myspec-gwt SKILL.md（worktree 创建）
- myspec-driven 自定义 OpenSpec Schema（6 artifact DAG）
- OpenSpec 环境已初始化（config.yaml、opsx 技能）

待实现：Go CLI 工具，负责将技能文件和 schema 分发到目标项目。

## Goals / Non-Goals

**Goals:**
- 构建 Go CLI 工具（`myspec install/update/list/uninstall/check/doctor`）
- 使用 `go:embed` 嵌入技能文件和 schema
- 文件复制安装（非符号链接），项目自包含，git 可追踪
- 注册表管理（`~/.config/myspec/registry.json`）
- 自动检测 OpenSpec CLI 并在缺失时提示安装
- 自动检测目标项目是否已初始化 OpenSpec，未初始化则自动执行 `openspec init`
- 版本管理通过 git tag，安装时检查 OpenSpec 版本兼容性

**Non-Goals:**
- 不自动安装或修改 OpenSpec 版本（只警告）
- 不管理 OpenSpec 的 opsx 技能（由 `openspec init` 负责）
- 不支持符号链接安装
- 不实现远程技能仓库
- 不实现技能版本回滚

## Decisions

### 决策 1：嵌入策略

嵌入技能文件 + 自定义 schema + config.yaml 模板到 Go 二进制中。

嵌入内容：
- `.claude/skills/myspec-br/SKILL.md`
- `.claude/skills/myspec-gwt/SKILL.md`
- `openspec/schemas/myspec-driven/`（schema.yaml + 6 个模板）
- `openspec-version.txt`（记录测试过的 OpenSpec 版本）

不嵌入 opsx 技能（由 `openspec init --tools claude` 负责）。

### 决策 2：版本管理

- 版本来源：git tag（如 `v1.0.0`）
- `openspec-version.txt` 记录测试过的 OpenSpec 版本
- `myspec install` 时比对系统 OpenSpec 版本与 `openspec-version.txt`
- 版本不匹配时警告但不阻止，提供修复命令
- 不自动安装或修改 OpenSpec

### 决策 3：OpenSpec 初始化

- `myspec install` 检测目标项目是否已 `openspec init`
- 未初始化则自动执行 `openspec init --tools claude`
- 已初始化则跳过

### 决策 4：注册表格式

```json
{
  "version": 1,
  "installed": {
    "/path/to/project": {
      "version": "v1.0.0",
      "installedAt": "2026-06-15T10:30:00Z",
      "skills": ["myspec-br", "myspec-gwt"],
      "schema": "myspec-driven"
    }
  }
}
```

### 决策 5：config.yaml 合并策略

- 目标项目无 `openspec/config.yaml` → 创建，写入 `schema: myspec-driven`
- 已有但无 schema 字段 → 添加 `schema: myspec-driven`
- 已有且有 schema 字段 → 更新为 `myspec-driven`
- 不覆盖 `context` 和 `rules` 字段

### 决策 6：CLI 命令

```bash
myspec install <project-path>     # 安装技能到项目
myspec update                     # 更新所有已安装项目
myspec update <project-path>      # 更新指定项目
myspec list                       # 列出已安装项目及版本
myspec uninstall <project-path>   # 从项目卸载
myspec check                      # 检查是否有新版本可用
myspec doctor                     # 诊断 OpenSpec 兼容性
```

### 决策 7：Go 项目结构

```
cmd/myspec/main.go          ← CLI 入口
internal/
  install/install.go        ← 安装逻辑
  update/update.go          ← 更新逻辑
  list/list.go              ← 列表逻辑
  uninstall/uninstall.go    ← 卸载逻辑
  check/check.go            ← 版本检查
  doctor/doctor.go          ← 诊断
  registry/registry.go      ← 注册表管理
  openspec/openspec.go      ← OpenSpec 检测与初始化
  embed/embed.go            ← go:embed 资源
embed/                      ← 嵌入资源目录
  skills/myspec-br/SKILL.md
  skills/myspec-gwt/SKILL.md
  schemas/myspec-driven/
openspec-version.txt
go.mod
```

## Risks / Trade-offs

- **[OpenSpec 版本不兼容]** → 版本检查警告，提供修复命令，不阻止安装
- **[目标项目已有 config.yaml]** → 合并策略只更新 schema 字段
- **[npm 全局安装权限问题]** → 不自动安装 OpenSpec，只提示用户手动安装
- **[注册表损坏]** → `myspec doctor` 检测并报告
- **[go:embed 路径错误]** → 使用 `all:` 前缀嵌入整个目录，避免路径问题
