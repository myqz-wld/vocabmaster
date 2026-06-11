# Plans 索引

> **范围**：终态 plan 文档。未终态 plan 留在当前环境配置的工作区；无更强契约时用 `<repo>/.refs/plans/`，`.refs/` 必须加入 `.gitignore`，不要放进本目录。
> **清理**：plan 到终态后，把最终文档和 plan 专属支持材料归档到 `ref/plans/`，更新本 INDEX，并清理工作区副本。

| Plan | 状态 | 完成日期 | 摘要 | 关联 changelog/review |
|---|---|---:|---|---|
| [src-layout-alignment-20260608.md](src-layout-alignment-20260608.md) | completed | 2026-06-08 | first-party source/data/tools 收敛到 `src/` 并更新工程入口 | [CHANGELOG_4](../changelogs/CHANGELOG_4.md) |
| [build-dir-migration-20260526.md](build-dir-migration-20260526.md) | completed | 2026-05-26 | vocabmaster (Go CLI) 编译产物从根 `./vocabmaster` 迁到 `./build/vocabmaster` | [CHANGELOG_1](../changelogs/CHANGELOG_1.md) |
