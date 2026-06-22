# Plans Index

Terminal plan documents live here. Draft or in-progress plans stay in the current environment's plan workspace; when no stronger convention exists, use `<repo>/.ref/plans/`. `.ref/` must stay ignored and must not contain terminal records.

When a plan reaches terminal state, archive the final document and plan-specific support material here, update this index, and remove workspace drafts.

## Naming

Existing historical records keep their current filenames. New final plans use `PLAN_X_<topic>.md`. Before creating one, run `ls ref/plans/`, set `X` to the maximum existing plan number in this directory plus 1, and do not guess. `<topic>` is short stable kebab-case and must not be vague like `update`, `fix`, or `misc`. Update this index in the same change.

| Plan | Status | Completed | Summary | Related Changelog/Review |
|---|---|---:|---|---|
| [src-layout-alignment-20260608.md](src-layout-alignment-20260608.md) | completed | 2026-06-08 | Consolidated first-party source/data/tools under `src/` and updated engineering entry points. | [CHANGELOG_4](../changelogs/CHANGELOG_4.md) |
| [build-dir-migration-20260526.md](build-dir-migration-20260526.md) | completed | 2026-05-26 | Moved the vocabmaster Go CLI build artifact from root `./vocabmaster` to `./build/vocabmaster`. | [CHANGELOG_1](../changelogs/CHANGELOG_1.md) |
