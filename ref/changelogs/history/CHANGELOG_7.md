# CHANGELOG_7: project foundation template alignment

## 概要

仓库级提示词资产和 `ref/` 索引按 project-engineering-foundation 模板重新组织，减少入口规则漂移。
同时把 `.refs/` 明确为未终态 plan/review 工作区，继续保持 `ref/` 为 tracked durable archive。

## 变更内容

### 根目录提示词资产

- `AGENTS.md` 补齐目录架构入口检查点，只保留当前入口的运行时 / 工具差异。
- `CLAUDE.md` 按模板拆分为仓库基础、基础目录架构、改动后必做、项目触发、项目约定和验证流程。
- 保留 VocabMaster 的 Go 1.24、`src/` / `build/`、LLM 离线可用和数据资产不变量。

### `ref/` 索引

- `ref/conventions/INDEX.md` 和 `ref/conventions/tally.md` 改用当前仓库路径，移除旧全局 prompt 路径和历史迁移说明。
- `ref/plans/INDEX.md` 改为终态 plan 模板列，保留现有两条 plan 的完成状态、日期和关联 changelog。
- `ref/reviews/INDEX.md` 补充未终态 review 草稿位置，保留 review 编号独立递增规则。
- `ref/changelogs/INDEX.md` 更新本条记录摘要。

### `.gitignore`

- 添加 `.refs/`，与 `build/`、`dist/` 一起保持本地工作区和构建产物不入 git。

## 备注

- `README.md` 仅做 check-only；本轮未改变用户可见命令、安装方式或项目结构说明。
- 本轮验证命令：`git diff --check`、Markdown 本地链接 / 路径检查、`make test`、`git status --short --branch`。
