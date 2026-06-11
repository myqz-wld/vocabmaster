# CHANGELOG_9: 入口资产去重 + README study 调度表边界修正

## 概要

prompt-asset 维护轮：三入口文件小幅去重，README study 调度表修正一处与代码不符的边界值。独立 review 裁决 0 MUST-FIX。

## 变更内容

### `README.md`

- §study 命令的智能调度 表格边界修正：代码 `src/internal/session/study.go` `calculateNewWords` 为 `>20 → 0 / >10 → 5 / 其余 → 10`，原表「10-20 → 5 / <10 → 10」把 dueCount=10 错归到 5 新词档；改为「11-20 → 5 / ≤10 → 10」与代码一致。
- 其余声明实测核对无误（词库数量 / 10 个子命令 / Makefile 安装参数 / LLM provider 顺序 / --data-dir），未改动。

### `CLAUDE.md`

- §项目特定约定 删 2 条与 §基础目录架构 重复的 bullet（build//dist/ 输出根、ref//.refs/ 边界）；保留 go.mod 导入路径和 LLM 离线兜底两条独有不变量。

### `AGENTS.md`

- §必读顺序 删第 3 条（成对资产审计，与 §入口特定操作说明 重复）；其独有的「协议语义不能漂移」从句并入后者。

## 备注

- 验证：死链 0；`git diff --check` 通过；独立 reviewer 用 study.go 实测确认表格修正方向正确。
