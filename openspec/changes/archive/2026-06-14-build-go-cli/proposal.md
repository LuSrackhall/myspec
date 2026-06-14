## Why

myspec 是一个 Claude Code 工作流管理工具，已完成技能定义（myspec-br、mysspec-gwt）和自定义 OpenSpec schema（myspec-driven），但缺乏跨项目分发机制。用户必须手动复制技能文件到目标项目，容易出错且不可追踪。需要一个 Go CLI 工具自动化这一过程。

## What Changes

- 新增 Go CLI 工具，支持 7 个子命令：`install`、`update`、`list`、`uninstall`、`check`、`doctor`
- 使用 `go:embed` 嵌入技能文件和 schema 到单一二进制
- 注册表管理（`~/.config/myspec/registry.json`）追踪已安装项目
- 自动检测 OpenSpec CLI 并在缺失时提示安装
- 自动检测目标项目 OpenSpec 初始化状态，未初始化则自动执行 `openspec init`
- 版本检查：安装时比对系统 OpenSpec 版本与嵌入版本

## Capabilities

### New Capabilities
- `cli-core`: Go CLI 入口、子命令路由、帮助信息
- `cli-install`: 安装逻辑（文件复制、OpenSpec 初始化、config.yaml 合并）
- `cli-registry`: 注册表管理（读写 `~/.config/myspec/registry.json`）
- `cli-openspec-bridge`: OpenSpec 检测、版本比对、自动初始化

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- 新增 Go 源码目录：`cmd/myspec/`、`internal/`
- 新增嵌入资源目录：`embed/`
- 新增版本文件：`openspec-version.txt`
- `.gitignore` 需添加 Go 构建产物（已有条目）
- 不影响现有技能文件和 schema
