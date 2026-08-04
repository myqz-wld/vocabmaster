# CHANGELOG_10: AGENTS.md 入口资产再收敛

## 概要

`AGENTS.md` 进一步收敛为只指向 `CLAUDE.md` 的入口说明，删除通用工具流程和 prompt-asset 维护细节，避免与共享 SSOT 双写漂移。

## 变更内容

### `AGENTS.md`

- 明确 `CLAUDE.md` 是 VocabMaster 仓库级共享 SSOT，`AGENTS.md` 只记录当前 agent 入口差异。
- 删除通用 `rg` / `apply_patch` / worktree / handoff / async 等运行时工具流程；这些不属于本仓库项目规则。
- 删除 prompt-asset 维护流程复述，保留“当前无额外入口差异”的边界说明。

## 验证

- `rg -n "apply_patch|rg|sleep|handoff|prompt-asset" AGENTS.md` 无匹配。
- `git diff --check -- AGENTS.md ref/changelogs/INDEX.md ref/changelogs/CHANGELOG_10.md` 通过。
