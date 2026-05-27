# CLAUDE.md

> 给 Claude Code 在本仓库工作时的硬性约定。本文件聚焦 **VocabMaster 专属** 的设计要点与改动流程；通用工程约定（输出语言、运行时、决策对抗、新项目工程地基）见 `~/.claude/CLAUDE.md`，本文件不重复。

## 仓库基础

- macOS 环境；语言 **Go ≥ 1.24**（gvm 管理；通用约定 §运行时 §Go：用项目对应版本）
- 入口：`main.go` / `cmd/` 各子命令
- 构建走 `Makefile`：`make build` / `make install` / `make clean`

## 项目定位

命令行背单词工具，基于 SM-2 间隔重复算法，支持英文（ECDICT 12,100+）与日文（JLPT 8,500+）。学习时可选调用本地 `claude` CLI 生成例句 / 润色释义；离线也能跑（不影响学习路径）。

## 构建 / 验证

```bash
make build         # go build → build/vocabmaster
make install       # 安装到 GOPATH/bin
make clean         # rm -rf build
go test ./...      # 单元测试
```

## src/build 目录约定（符合通用 §新项目工程地基 §src/build 标准目录结构）

- **源码** Go 项目按 `cmd/` `internal/` `tools/` 等社区惯例（不强行收敛到 `src/` —— Go 生态以 `cmd/` `internal/` `pkg/` 为约定俗成；属于"按工具链默认惯例保留原状"例外）
- **build 产物** `build/`（Makefile `go build -o build/<binary>`）
- **数据资产** `data/`（词库 JSON）— 入 git，运行时只读

## 改动后必做（最低操作指南）

> 详细约定见通用 §新项目工程地基（user / 应用 SDK 注入）。

1. **改用户可见行为 / 文件结构 / 启动方式** → 改对应章节 `README.md`；纯 bug 修复 / 内部重构不动 README
2. **写 changelog 或 review 二选一**（必做）：
   - 功能变更 / 行为修改 / 命令新增 / 依赖升级 → `ref/changelogs/CHANGELOG_X.md`（X 递增）+ 同步 `ref/changelogs/INDEX.md`
   - Debug / 性能 / 安全 review → `ref/reviews/REVIEW_X.md`（X 递增）+ 同步 `ref/reviews/INDEX.md`
3. **改功能前先读** `ls ref/conventions/ ref/changelogs/ ref/reviews/` + 浏览相关条目
4. **数据资产改动**（`data/` 词库 JSON）按数据迁移姿势处理，不入 changelog 主线（视情况落 review）
