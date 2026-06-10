# Reviews 索引

> 周期性 / 触发性的 debug、code review、性能 audit、安全审查报告。功能变更去 [`changelogs/`](../changelogs/INDEX.md)。

## 命名

`REVIEW_X.md`（X 在 `ref/reviews/` 内独立递增整数）。新建前 `ls ref/reviews/` 找最大 X；关联的 changelog 写在索引列里，不共用序号。

## 单文件结构

- 触发场景（用户主动 / 周期性 / 大重构前 ...）
- 方法（双对抗 Agent 配对、范围、工具）
- 三态裁决清单（✅ / ❌ / ❓）+ 证据（文件:行号 + 代码片段）
- 修复条目（按严重度）
- 关联 changelog（本轮修复落地的 CHANGELOG 编号）

## 索引

| 文件 | 主题 | 严重度分布 | 关联 changelog |
|------|------|-----------|----------------|
| [REVIEW_1.md](REVIEW_1.md) | install vm environment + project deep review | R1: 0 P0, 0 P1 after rebuttal, 4 P2, 3 P3; R2: 2 P3 fixed | [3](../changelogs/CHANGELOG_3.md) |
