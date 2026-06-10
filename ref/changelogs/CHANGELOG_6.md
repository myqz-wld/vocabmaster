# CHANGELOG_6: LLM enhancement provider fallback

## Summary

LLM enhancement now tries local providers in the explicit order `codex -> claude -> fail`. Codex CLI is preferred when available; Claude Code remains the fallback; if neither provider can produce valid enhancement JSON, the caller receives a failure and the study flow keeps using base word data.

## Changes

- Added an ordered LLM provider layer in `src/internal/llm`, with Codex CLI first and Claude Code second.
- Codex runs through `codex exec` in non-interactive read-only mode and writes the final assistant message to a temporary output file for parsing.
- Provider failures, timeouts, empty responses, and invalid enhancement JSON now fall through to the next provider before failing the request.
- Updated `vm generate` and README messaging to describe the new `codex -> claude` order.
- Added unit coverage for provider order, fallback, all-provider failure, and availability detection.

## Verification

- `go test ./src/internal/llm -v`
- `make test`
