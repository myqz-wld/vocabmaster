# Plans 索引

> **范围**：完整的设计 plan(走 user CLAUDE.md §复杂 plan 流程归档)。简单约定升级 / 候选追踪见 [`conventions/`](../conventions/INDEX.md);功能变更见 [`changelogs/`](../changelogs/INDEX.md);review 见 [`reviews/`](../reviews/INDEX.md)。
>
> **canonical**: 4 列对齐应用约定 archive_plan tool 4 列 canonical(`| 文件 | 状态 | 关联 changelog | 概要 |`)— vocabmaster 项目目前 G-manual 归档,本 INDEX 由 lead 手工维护。

| 文件 | 状态 | 关联 changelog | 概要 |
|------|------|--------------|------|
| [src-layout-alignment-20260608.md](src-layout-alignment-20260608.md) | completed | [4](../changelogs/CHANGELOG_4.md) | first-party source/data/tools 收敛到 `src/` 并更新工程入口 |
| [build-dir-migration-20260526.md](build-dir-migration-20260526.md) | completed | [1](../changelogs/CHANGELOG_1.md) | vocabmaster (Go CLI) 编译产物从根 `./vocabmaster` 迁到 `./build/vocabmaster`,对齐 cross-language `build/` canonical |
