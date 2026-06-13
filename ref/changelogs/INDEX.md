# Changelog 索引

> **范围**：功能变更（新功能 / 行为修改 / API / 依赖升级）。Debug / 性能 / 安全 review 见 [`reviews/`](../reviews/INDEX.md)。

| 文件 | 概要（≤80 字） |
|------|------|
| [CHANGELOG_9.md](CHANGELOG_9.md) | 入口资产去重；README study 边界修正；补 `scripts/` 目录规则 |
| [CHANGELOG_8.md](CHANGELOG_8.md) | foundation 二轮对齐：补 review 过期 + 500 行护栏节，落地 expiry 脚本，删 AGENTS 重复节 |
| [CHANGELOG_7.md](CHANGELOG_7.md) | 对齐 foundation 模板组织并补 `.refs/` 非终态工作区规则 |
| [CHANGELOG_6.md](CHANGELOG_6.md) | LLM 增强 provider 顺序改为 `codex -> claude -> fail` |
| [CHANGELOG_5.md](CHANGELOG_5.md) | review_history 跨 offset 日期统计改为 UTC 时间轴比较 |
| [CHANGELOG_4.md](CHANGELOG_4.md) | 将 first-party source/data/tools 收敛到 `src/` 并补 `AGENTS.md` |
| [CHANGELOG_3.md](CHANGELOG_3.md) | make install 写入 shell 配置，让 `vm` 在新终端中直接可用 |
| [CHANGELOG_2.md](CHANGELOG_2.md) | make install 默认安装到 `$GOPATH/bin`，并创建 `vm` 短命令入口 |
| [CHANGELOG_1.md](CHANGELOG_1.md) | build-dir-migration: vocabmaster (Go CLI) 编译产物从根 `./vocabmaster` 迁到 `./build/vocabmaster`,对齐 cross-language `build/` canonical |
