# Plans Index

## Scope

Routing index for final plan documents. Only final plan documents enter `ref/plans/`. Keep non-final plans in the current environment's plan workspace or `<repo>/.ref/plans/`. Archive durable support materials meant for future agents somewhere under `ref/` and link them back to the plan.

This root index defines routing and bucket policy only. Per-record plan rows live only in the bucket `INDEX.md` files.

## Naming

New final plans use `ref/plans/<bucket>/PLAN_X_<topic>.md`. Before creating one, scan `ref/plans/*/PLAN_*.md`, set `X` to the maximum existing plan number plus 1, and do not guess. `<topic>` is short stable kebab-case and must not be vague like `update`, `fix`, or `misc`. Legacy records keep their original filenames when rebucketed.

## File Structure

- Goal
- Context / constraints
- Task breakdown
- Validation
- Completed date as `Completed At: <YYYY-MM-DD>` or frontmatter `completed_at`
- Final status / handoff

## Buckets

Buckets are mutually exclusive. Store each plan in exactly one bucket based on `Completed At` or `completed_at`.

| Bucket | Date Range | Directory |
|---|---|---|
| Recent 3 days | Completion date is within the last 3 days, inclusive | `ref/plans/recent-3-days/` |
| Recent week | Older than 3 days and within the last 7 days, inclusive | `ref/plans/recent-week/` |
| Recent month | Older than 7 days and within the last 30 days, inclusive | `ref/plans/recent-month/` |
| History | Older than 30 days, or missing a parseable date | `ref/plans/history/` |

## Rebucket Rules

On every new or edited final plan:

1. Scan all `ref/plans/*/PLAN_*.md` files.
2. Recompute each file's bucket from `Completed At` or `completed_at`.
3. Move files that no longer belong in their current bucket.
4. Update this routing index only if bucket policy changes.
5. Update every affected bucket `INDEX.md` while preserving its table format.

For full recent-week context, read `recent-3-days/` plus `recent-week/`. For full recent-month context, read `recent-3-days/`, `recent-week/`, and `recent-month/`. For full history, read all four buckets.

## Bucket Indexes

- `ref/plans/recent-3-days/INDEX.md`
- `ref/plans/recent-week/INDEX.md`
- `ref/plans/recent-month/INDEX.md`
- `ref/plans/history/INDEX.md`
