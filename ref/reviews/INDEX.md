# Reviews Index

Debug, code review, performance audit, and security review reports. Draft review notes stay under `<repo>/.ref/reviews/`; terminal reviews are archived here. Functional changes belong in [Changelog Index](../changelogs/INDEX.md).

## Naming

Existing historical records keep their current filenames. New final reviews use `REVIEW_X_<topic>.md`. Before creating one, run `ls ref/reviews/`, set `X` to the maximum existing review number in this directory plus 1, and do not guess. `<topic>` is short stable kebab-case and must not be vague like `update`, `fix`, or `misc`. Related changelogs are recorded in the index table and do not share the review number. Update this index in the same change.

## Record Format

- Trigger: user request, periodic review, pre-refactor gate, or incident.
- Method: reviewer pairing, scope, and tools.
- Decision list: `✅ / ❌ / ❓` with evidence such as `file:line` and short snippets.
- Fixes grouped by severity.
- Related changelog for fixes landed from the review.

## Index

| File | Topic | Severity Distribution | Related Changelog |
|---|---|---|---|
| [REVIEW_1.md](REVIEW_1.md) | install vm environment + project deep review | R1: 0 P0, 0 P1 after rebuttal, 4 P2, 3 P3; R2: 2 P3 fixed | [3](../changelogs/CHANGELOG_3.md) |
