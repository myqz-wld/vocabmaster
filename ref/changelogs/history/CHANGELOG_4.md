# CHANGELOG_4: src layout alignment

## Summary

Repository organization now follows the project-engineering-foundation layout: first-party source, data, and tools live under `src/`; root-level files are limited to metadata, build entry points, AI-coding entry points, and reference records.

## Changes

- Moved `main.go`, `cmd/`, `internal/`, `data/`, and `tools/` under `src/`.
- Updated Go imports to `github.com/vocabmaster/vocabmaster/src/...`.
- Updated `Makefile` to build `./src` and call `src/tools/install.sh` / `src/tools/install_test.sh`.
- Added `AGENTS.md` as the companion project agent entry and updated `CLAUDE.md` to make `src/` the source SSOT.
- Updated `README.md` with the current project structure.
- Updated `src/tools/extract_words.py` default output to `src/data/`.
- Added prompt-asset local backup/inventory cache under `.prompt-asset-improver/local/` and ignored it in `.gitignore`.
- Ignored `.deep-review-cache/` for Agent Deck review cache files.
- Fixed active `ref/*/INDEX.md` links/descriptions to use `ref/changelogs/`.

## Verification

- `python3 -m py_compile src/tools/extract_words.py`
- `bash -n src/tools/install.sh && bash -n src/tools/install_test.sh`
- `make test`
- `make build`
- `./build/vocabmaster --help`
- Isolated `make install BINDIR=<tmp>/bin SHELL_RC=<tmp>/.zshrc`
- `<tmp>/bin/vm --help`
- Isolated `make uninstall BINDIR=<tmp>/bin SHELL_RC=<tmp>/.zshrc`
- Real `make install`
- `zsh -lc 'source ~/.zshrc; vm --help'`
- `rg -n '\.\./changelog/|changelog/' ref/conventions/INDEX.md ref/plans/INDEX.md ref/reviews/INDEX.md ref/changelogs/INDEX.md` returns no matches
