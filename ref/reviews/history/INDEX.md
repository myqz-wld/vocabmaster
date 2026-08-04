# History Reviews

## Scope

This bucket contains reviews older than 30 days or with no parseable `reviewed_at`. Remove rows for files moved during rebucketing.

| Bucket | Date Range |
|---|---|
| `recent-3-days` | Within the last 3 days, inclusive |
| `recent-week` | Older than 3 days and within the last 7 days, inclusive |
| `recent-month` | Older than 7 days and within the last 30 days, inclusive |
| `history` | Older than 30 days, or missing a parseable date |

## Index Table

| reviewed_at | File | Topic | Severity Distribution |
|---|---|---|---|
| unknown | `REVIEW_2_codex-cli-flag.md` | Codex CLI flag compatibility | 1 P2 fixed |
| unknown | `REVIEW_1.md` | Install environment and project deep review | 4 P2, 5 P3 fixed/accepted |
