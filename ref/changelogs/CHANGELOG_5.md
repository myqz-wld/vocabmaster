# CHANGELOG_5: review history offset-safe date counts

## Summary

Review count and streak statistics now compare `review_history.reviewed_at` on a normalized UTC timeline instead of lexicographically comparing RFC3339 text with offsets.

## Changes

- `store.GetReviewCountOnDate` now uses SQLite `unixepoch(...)` for the `reviewed_at` window comparison, so existing history rows with different RFC3339 offsets are counted by instant.
- Added regression coverage for cross-offset history rows where the stored text date differs from the target local calendar day.

## Verification

- `go test ./src/internal/store -v`
- `make test`
