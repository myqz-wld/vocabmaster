---
changelog_id: 13
changed_at: 2026-08-04
---

# CHANGELOG_13_foundation-llm-adapters: Align foundation and configure LLM adapters

## Summary

Aligned the repository with the current project-engineering foundation and added explicit Claude, Codex, and Grok selection for vocabulary enrichment. Users can now choose an adapter, model, and supported thinking effort while retaining offline fallback and cache behavior.

## Changes

### Engineering foundation

- Added mutually exclusive time buckets and current routing indexes under `ref/changelogs/`, `ref/reviews/`, and `ref/plans/`; moved legacy terminal records into `history/` without rewriting their content.
- Updated `CLAUDE.md`, `AGENTS.md`, and `UI_COPY_LANGUAGE.md` with current template rules while preserving VocabMaster-specific invariants.
- Replaced the plan-only archive reminder with the current `.ref` reminder and refreshed the file-level review-expiry helper.
- Moved installation, installer tests, and vocabulary extraction from `src/tools/` to `scripts/`, updating all live paths.

### LLM enrichment

- Split the former 476-line LLM implementation into client, adapter, prompt, and response modules.
- Added persistent `--llm-adapter`, `--llm-model`, and `--llm-thinking` flags for `study`, `learn`, `review`, and `generate` flows.
- Added Grok CLI support and extended default automatic fallback to `codex -> claude -> grok`.
- Added adapter-specific effort validation and native model/effort argument forwarding for all three CLIs.
- Preserved cache-first behavior; `generate --force` remains the explicit regeneration path.

### User documentation and validation fixes

- Rewrote README as a concise user-facing guide centered on installation, common commands, and LLM configuration.
- Fixed the installer's missing-metadata error interpolation so its existing negative test returns the intended nonzero status.

## Validation

- `git diff --check`
- `~/.local/bin/mise exec -- go test ./...`
- `~/.local/bin/mise exec -- go test -race ./src/internal/llm`
- `make test`
- `make build`
- Isolated `make install` / installed `vm --version` / `make uninstall`
- `bash -n` for all maintained shell scripts
- Python import/path assertion for `scripts/extract_words.py`
- Claude, Codex, and Grok local `--help` argument-contract probes without model calls
- CLI help and invalid configuration smoke tests
- `bash scripts/file-level-review-expiry.sh`
- Source file-size scan; all non-generated source files remain below 500 LOC

## Do Not Split Protection

None.

## Notes

Model availability and supported reasoning effort remain adapter/model dependent. VocabMaster validates its documented adapter-level values first, then lets the selected CLI report model-specific incompatibility.
