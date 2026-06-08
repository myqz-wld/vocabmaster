# AGENTS.md

> 本文件作为仓库级 agent 指令加载。共享仓库规则写在 `CLAUDE.md`：仓库基础、改动后要求、项目约定和验证流程。本文件只记录当前入口的运行时 / 工具差异，避免与 `CLAUDE.md` 漂移。

## 必读顺序

1. 先读 `CLAUDE.md`，再遵守其中的项目结构、验证流程和改动后要求。
2. 涉及 SDK session、MCP tool、skill 或 prompt 资产时，遵守当前用户环境配置的契约；没有明确契约时，不要发明工具流程。
3. 成对 prompt 资产必须同时审计。运行时机制不同时，工具措辞可以不同，但协议语义不能漂移。

## 入口特定操作说明

- 默认用 `rg` 搜索代码，用 `apply_patch` 手工编辑；不要通过 shell 重定向或一次性脚本写项目文件。
- 使用当前环境提供的 worktree 或 handoff 工具。普通 shell 命令使用 `git -C <worktree>` 或绝对路径。
- 本 SDK 是 turn-based：如果环境提供跨会话消息或异步协作工具，发送消息后报告状态并结束当前 turn，等待回复。没有明确契约时，不要用 `sleep` 或轮询模拟阻塞等待。
- 编辑长期 prompt 资产前，运行 prompt-asset 维护检查；存在成对 counterpart 时同时审计。

## 项目特定入口差异

当前无额外入口差异；共享规则以 `CLAUDE.md` 为准。
