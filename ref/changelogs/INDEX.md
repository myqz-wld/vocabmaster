# Changelog Index

Functional changes, including new features, behavior changes, APIs, dependencies, or structure changes. Debug, performance, security, and review-driven fixes belong in [Reviews Index](../reviews/INDEX.md).

## Naming

Existing historical records keep their current filenames. New final changelogs use `CHANGELOG_X_<topic>.md`. Before creating one, run `ls ref/changelogs/`, set `X` to the maximum existing changelog number in this directory plus 1, and do not guess. `<topic>` is short stable kebab-case and must not be vague like `update`, `fix`, or `misc`. Update this index in the same change.

| File | Summary |
|---|---|
| [CHANGELOG_11.md](CHANGELOG_11.md) | Tightened foundation wording and LLM enrichment prompt constraints. |
| [CHANGELOG_10.md](CHANGELOG_10.md) | Further narrowed AGENTS.md to read-only CLAUDE.md and entry-point differences. |
| [CHANGELOG_9.md](CHANGELOG_9.md) | Deduplicated entry assets, corrected README study boundaries, and added `scripts/` rules. |
| [CHANGELOG_8.md](CHANGELOG_8.md) | Aligned foundation again with review expiry, 500-line guards, expiry script, and AGENTS cleanup. |
| [CHANGELOG_7.md](CHANGELOG_7.md) | Aligned foundation template organization and added `.refs/` non-final workspace rules. |
| [CHANGELOG_6.md](CHANGELOG_6.md) | Changed LLM enhancement provider order to `codex -> claude -> fail`. |
| [CHANGELOG_5.md](CHANGELOG_5.md) | Changed review_history cross-offset date stats to compare on a UTC timeline. |
| [CHANGELOG_4.md](CHANGELOG_4.md) | Consolidated first-party source, data, and tools under `src/` and added `AGENTS.md`. |
| [CHANGELOG_3.md](CHANGELOG_3.md) | make install now writes shell config so `vm` works directly in new terminals. |
| [CHANGELOG_2.md](CHANGELOG_2.md) | make install now defaults to `$GOPATH/bin` and creates the `vm` short command. |
| [CHANGELOG_1.md](CHANGELOG_1.md) | build-dir-migration: moved the Go CLI artifact from `./vocabmaster` to `./build/vocabmaster`. |
