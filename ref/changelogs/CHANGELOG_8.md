# CHANGELOG_8: foundation 模板二轮对齐（review 过期 + 文件大小护栏）

## 概要

按 project-engineering-foundation 模板补齐 CHANGELOG_7 一轮对齐后仍缺的规则节，并落地配套维护脚本；不改变运行时行为。

## 变更内容

### `CLAUDE.md`

- 新增「Review 过期与最小复审范围」节：unreviewed ∪ expired ∪ scope_unknown 最小复审范围 + 4 条过期判定（净改动 / commit 数 / 90 天 / frontmatter expired）。
- 新增「文件大小护栏（500 行）」节：拆分优先级 + 例外清单（含 `src/data/` 词库 JSON）。
- 「改动后必做」第 2 条补编号规则：`X` 取当前最大值 +1，用 `ls` 确认；INDEX 摘要 ≤ 80 字。
- 头部 SSOT 描述同步覆盖新增两节。

### `AGENTS.md`

- 删除「目录架构入口规则」节：内容与 `CLAUDE.md` 基础目录架构重复，违反「共享规则不要复制到两个文件」约定。

### `scripts/`

- 新增 `scripts/file-level-review-expiry.sh`（来自 foundation skill），review 过期检查脱离 skill 独立可跑。

### `README.md`

- 「项目结构」节补 `scripts/` 条目。

## 备注

- 本轮验证：`make test` 通过；`git diff --check` 无问题。
