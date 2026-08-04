---
changelog_id: 15
changed_at: 2026-08-04
---

# CHANGELOG_15_reasoning-levels: Align adapter reasoning levels

## Summary

Aligned VocabMaster's adapter-level thinking validation with the configured Codex and Grok reasoning tiers.

## Changes

- Codex now accepts `low`, `medium`, `high`, `xhigh`, and `max`; `minimal` is no longer accepted.
- Grok now accepts `low`, `medium`, `high`, and `xhigh`.
- Claude remains unchanged at `low`, `medium`, `high`, `xhigh`, and `max`.
- Updated validation tests and the concise user README.

## Validation

- `make test`
- `~/.local/bin/mise exec -- go test -race ./src/internal/llm`
- `~/.local/bin/mise exec -- go vet ./...`
- `git diff --check`

## Do Not Split Protection

None.

## Notes

Individual models can still reject an adapter-level value they do not support.
