# Recent 3 Days Reviews

## Scope

This bucket contains only reviews whose `reviewed_at` is within the last 3 days, inclusive. Remove rows for files moved during rebucketing.

| Bucket | Date Range |
|---|---|
| `recent-3-days` | Within the last 3 days, inclusive |
| `recent-week` | Older than 3 days and within the last 7 days, inclusive |
| `recent-month` | Older than 7 days and within the last 30 days, inclusive |
| `history` | Older than 30 days, or missing a parseable date |

## Index Table

| reviewed_at | File | Topic | Severity Distribution |
|---|---|---|---|
