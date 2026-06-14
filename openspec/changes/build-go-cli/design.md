## Context

myspec 仓库已完成：myspec-br 技能（脑暴编排器）、mysspec-gwt 技能（worktree 创建）、myspec-driven 自定义 OpenSpec schema（6 artifact DAG）、OpenSpec 环境初始化。需要一个 Go CLI 工具将这些组件分发到目标项目。

## Goals / Non-Goals

**Goals:**
- 单二进制分发，`go build` 即可使用
- 文件复制安装，目标项目自包含，git 可追踪
- 自动处理 OpenSpec 依赖检测和初始化
- 注册表追踪已安装项目，支持 update/list/uninstall/check

**Non-Goals:**
- 不管理 OpenSpec 版本（只警告）
- 不支持远程技能仓库
- 不实现技能版本回滚

## Decisions

### 决策 1：CLI 框架选择

**选择：** 标准库 `flag` + 手动子命令路由

**替代方案：**
- Cobra：功能强大但依赖重，myspec 只有 7 个子命令
- urfave/cli：类似 Cobra，同样过重

**理由：** 7 个子命令用标准库足够。`os.Args[1]` 路由到对应处理函数，每个子命令用 `flag.NewFlagSet` 解析参数。保持零外部依赖。

### 决策 2：go:embed 目录结构

```
embed/
├── skills/
│   ├── myspec-br/SKILL.md
│   └── myspec-gwt/SKILL.md
└── schemas/
    └── myspec-driven/
        ├── schema.yaml
        └── templates/*.md
```

使用 `//go:embed all:embed` 嵌入整个目录。`all:` 前缀确保包含隐藏文件。

### 决策 3：OpenSpec 检测策略

```
1. exec.LookPath("openspec") → 检测 CLI 是否存在
2. exec.Command("openspec", "--version") → 获取版本
3. 比对 openspec-version.txt 中的嵌入版本
4. 不匹配 → 警告 + 提供 npm install 命令
5. 不存在 → 报错 + 提供安装命令
```

不自动安装。不自动修改版本。

### 决策 4：config.yaml 合并

读取目标项目的 `openspec/config.yaml`（YAML 解析），更新 `schema` 字段为 `myspec-driven`，保留 `context` 和 `rules`。使用 `gopkg.in/yaml.v3` 解析（唯一的外部依赖）。

### 决策 5：注册表并发安全

使用 `os.OpenFile` 的 `O_CREATE|O_RDWR` + 文件锁（`syscall.Flock`）确保并发安全。单人开发场景下不太可能并发，但作为正确性保障。

## Risks / Trade-offs

- **[yaml.v3 依赖]** → 唯一外部依赖，可考虑用 JSON 替代 YAML 配置来消除
- **[OpenSpec 版本不兼容]** → 警告但不阻止，用户自行决定
- **[目标项目已有 skills]** → install 覆盖，用户可通过 registry 追踪
- **[go:embed 路径]** → 使用 `all:` 前缀，测试验证路径正确性
