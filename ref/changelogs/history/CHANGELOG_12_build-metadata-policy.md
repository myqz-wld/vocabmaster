# CHANGELOG_12_build-metadata-policy: Add installed CLI build metadata

## Summary

Added build metadata generation, installation, and local freshness checks for the VocabMaster CLI. `make build` now writes git-derived metadata, `make install` ships that metadata with the installed binary, and the CLI exposes human-readable and machine-checkable installed-version status commands.

## Changes

### `CLAUDE.md`

- Identified `make install` as the active installable CLI packaging surface.
- Required installed CLI metadata to include app/package name, semantic version when available, full and short git commit, branch when available, dirty flag when determinable, and build timestamp.
- Required a human-readable version/status entry and a machine-checkable freshness command or equivalent.
- Required freshness checks to avoid remote fetches and to report missing metadata distinctly from commit mismatches.

### Build and install tooling

- Added `scripts/write-build-info.go` to generate `build/build-info.json` with package name, version when available, full and short commit, branch when available, dirty flag, and build timestamp.
- Updated `make build` to generate build metadata before compiling `build/vocabmaster`.
- Updated `src/tools/install.sh` to require build metadata during install, copy it to `vocabmaster.build-info.json`, and remove it during uninstall.
- Updated `src/tools/install_test.sh` to cover metadata install, uninstall, and missing-metadata rejection.

### CLI status commands

- Added `vm --version` / `vocabmaster --version` and `vm version` / `vocabmaster version` for human-readable installed metadata and local checkout comparison.
- Added `vm --check-installed` / `vocabmaster --check-installed` and `vm check-installed` / `vocabmaster check-installed` for machine-checkable freshness.
- The check exits `0` when installed metadata matches the current checkout commit, `1` for commit mismatch, and `2` when metadata or local comparison state is missing/unavailable.
- Freshness checks use only local git state and optionally report local `origin/main`; they do not fetch remotes.

### `README.md`

- Documented installed version checks, exit codes, and the no-fetch behavior.
- Added `version` and `check-installed` to the command table.

## Validation

- `make build`
- `./build/vocabmaster --version`
- Isolated install into a temporary `BINDIR`
- Installed `vm --version`
- Installed `vm --check-installed`
- `make test`
- `git diff --check`

## Notes

The generated `build/build-info.json` remains a local ignored build artifact.
