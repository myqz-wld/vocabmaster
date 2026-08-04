---
changelog_id: 14
changed_at: 2026-08-04
---

# CHANGELOG_14_persistent-llm-config-packaging: Persist LLM defaults and package releases

## Summary

Added durable machine-local LLM defaults and a reproducible platform package command for release and installation closeout.

## Changes

### Persistent LLM configuration

- Added `vm config set-llm` and `vm config show`.
- Stored adapter, model, and thinking defaults in `<data-dir>/config.json` with owner-only file permissions.
- Kept `--llm-adapter`, `--llm-model`, and `--llm-thinking` as per-command overrides.
- Clear saved model and thinking values when a command temporarily switches to another adapter unless those values are also explicitly overridden.

### Packaging

- Added `make pack` and `scripts/package.sh`.
- Packages the CLI, `vm` symlink, build metadata, and concise README into `dist/vocabmaster-<release>-<platform>.tar.gz`.
- Generates an adjacent SHA-256 checksum and removes temporary staging files automatically.

## Validation

- `make test`
- `~/.local/bin/mise exec -- go test -race ./...`
- `~/.local/bin/mise exec -- go vet ./...`
- `make pack`
- Archive content inspection and SHA-256 verification
- Temporary-data-directory `config set-llm` / `config show` round trip
- Saved-default and per-command adapter/thinking override smoke tests
- `git diff --check`
- Shell syntax checks

## Do Not Split Protection

None.

## Notes

When no version tag is available, package names use the 12-character build commit. Final release packaging must run from a clean committed checkout so build metadata records `dirty: false`.
