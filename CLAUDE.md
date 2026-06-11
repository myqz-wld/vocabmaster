# CLAUDE.md

> 本文件是 VocabMaster 仓库级共享 SSOT，记录仓库基础、基础目录架构、改动后要求、plan/review 文档生命周期、review 过期规则、文件大小护栏、项目特定触发、项目特定约定和验证流程。
> `AGENTS.md` 是 companion entry，只记录当前 agent 入口的运行时 / 工具差异；共享规则不要复制到两个文件。
> 最小工程流程以本文件为准，额外工程或 review skill 只作为增强层。

## 仓库基础

- OS：macOS。
- 语言版本：**Go ≥ 1.24**（gvm 管理；通用约定 §运行时 §Go：用项目对应版本）。
- 入口：`src/main.go` / `src/cmd/` 各子命令。
- 构建入口：`Makefile`。

## 项目定位

命令行背单词工具，基于 SM-2 间隔重复算法，支持英文（ECDICT 12,100+）与日文（JLPT 8,500+）。学习时可选调用本地 Codex CLI / Claude Code 生成例句 / 润色释义；离线也能跑（不影响学习路径）。

## 基础目录架构

创建或维护仓库时按这份结构落位；除非项目已有更强契约，不要为同类文件另建平行目录：

- `CLAUDE.md`：共享项目 SSOT，记录仓库基础、目录架构、改动后必做、plan/review 生命周期、项目特定触发、项目特定约定和验证流程。
- `AGENTS.md`：入口 / 工具差异，只引用并遵守 `CLAUDE.md` 的共享规则。
- `README.md`：面向用户和维护者的安装、使用、验证和结构说明。
- `src/`：first-party Go 源码、内置词库数据和维护脚本；`src/data/` 词库 JSON 入 git，运行时只读。
- `build/`：本地构建产物；`Makefile` 用 `go build -o build/vocabmaster ./src`。
- `dist/`：保留为可选打包输出根；当前没有 active 产物。
- `ref/changelogs/INDEX.md`：终态 changelog 索引。
- `ref/reviews/INDEX.md`：终态 review 索引；终态 review 文件放 `ref/reviews/REVIEW_X.md`。
- `ref/plans/INDEX.md`：终态 plan 索引；终态 plan 文件放 `ref/plans/`。
- `ref/conventions/INDEX.md`：已升级项目约定索引；约定正文用 `ref/conventions/<X>-<topic>.md`。
- `ref/conventions/tally.md`：重复反馈 / 重复 agent 踩坑计数入口。
- `.refs/`：必须加入 `.gitignore`；只放未终态 plan/review 工作副本，不放终态记录。

## 改动后必做

先执行这几条最低规则，再按项目特定触发补充：

1. 改用户可见行为、文件结构、启动方式、端口、依赖或验证步骤 → 更新 `README.md` 对应章节；纯 bug 修复或内部重构不动 README。
2. 每个有意义的功能 / 行为 / 命令 / 依赖 / 结构变化 → 写 `ref/changelogs/CHANGELOG_X.md` 并更新 `ref/changelogs/INDEX.md`；debug / 性能 / 安全 / review-driven fix → 写 `ref/reviews/REVIEW_X.md` 并更新 `ref/reviews/INDEX.md`。编号 `X` 取当前最大值的下一个整数，用 `ls` 确认，不要猜；INDEX 摘要 ≤ 80 字或一句简短英文。
3. 未终态 plan/review 留在当前环境的工作区；无更强契约时用 `<repo>/.refs/`。终态收口时把最终 plan 归档到 `ref/plans/`、最终 review 归档到 `ref/reviews/REVIEW_X.md`，更新对应 INDEX，并清理工作区副本。
4. 反复用户反馈或重复 agent 踩坑先记入 `ref/conventions/tally.md`；`count >= 3` 后走本仓库 review 流程，升级为 `ref/conventions/<X>-<topic>.md` 并更新 `ref/conventions/INDEX.md`。
5. 改功能前先读 `ls ref/conventions ref/changelogs ref/plans ref/reviews` 并浏览相关条目。

## 项目特定触发

- 改 `src/data/` 词库 JSON：按数据更新流程处理；仅数据刷新不入 changelog 主线，视质量风险落 review。
- 改 CLI 命令、安装方式、构建输出路径或 LLM provider 顺序：同步更新 `README.md` 对应章节，并写 changelog。

## 项目特定约定（设计要点速查）

> 动态升级走 `ref/conventions/<X>-<topic>.md`；本节只保留必须在入口可见的项目不变量。

- `go.mod` / `go.sum` 保持在仓库根目录；first-party Go 包在 `src/` 下，导入路径使用 `github.com/vocabmaster/vocabmaster/src/...`。
- LLM 增强是可选本地能力；Codex CLI / Claude Code 不可用或调用失败时，学习路径必须继续使用内置基础数据。

## Review 过期与最小复审范围

准备下一次 review 时按本节确定最小复审范围；`ref/reviews/` 是会过期的覆盖记录，不是永久豁免。

下一次 review 的最小范围：

```text
unreviewed files ∪ expired reviewed files ∪ scope_unknown files
```

自最近一次覆盖该文件的 REVIEW 基线以来，满足任一条件即过期：

- 净改动 ≥ `min(200 行, 当前 LOC 的 30%)`。
- 不同 commit 数 ≥ 3。
- 距今 ≥ 90 天且文件至少改过一次。
- REVIEW frontmatter 标记 `expired: true`。

准备 review 时在仓库根目录运行 `bash scripts/file-level-review-expiry.sh`；脚本缺失时按上述条件用 `git log` 手工判定。

## 文件大小护栏（500 行）

任何源码文件超过 500 LOC，提交前必须先尝试拆分；生成代码、lockfile、快照、migration、fixture（含 `src/data/` 词库 JSON）除外。

拆分优先级：

1. 抽出模块级纯函数 / 类型 / 常量。
2. 目录化为同目录子模块并保持 import 路径。
3. 仅在 plan/review 之后才用 facade + 共享上下文拆类。

确实不可拆分时，在相关 changelog 的 "do not split" 保护清单中记录文件和具体原因。

## 验证流程

```bash
make build         # go build -> build/vocabmaster
make test          # go test ./... + installer shell tests
make install       # 安装到 GOPATH/bin；改安装流程时必须实测
make clean         # rm -rf build
```

改用户可见 CLI 行为、安装流程、构建路径或 LLM provider 顺序后，至少运行 `make test`；安装相关改动还要实测 `make install` 或等价隔离安装命令。

## 部署 / 打包

当前无单独部署流程；本地安装和发布前置检查以 `make install`、`make build`、`make test` 为准。
