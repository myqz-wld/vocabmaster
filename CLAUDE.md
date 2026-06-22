# CLAUDE.md

> Shared repository workflow for paired AI-coding entries.
> Put runtime or tool differences in `AGENTS.md` to avoid drift.
> This file defines the minimum engineering flow for VocabMaster. Additional engineering or review skills are enhancement layers.

## Repository Baseline

- OS: macOS.
- Language version: **Go >= 1.24** (managed by gvm; general convention § Runtime § Go: use the project-specific version).
- Entry points: `src/main.go` / subcommands under `src/cmd/`.
- Build entry point: `Makefile`.

## Project Purpose

A command-line vocabulary memorization tool based on the SM-2 spaced repetition algorithm. It supports English (ECDICT 12,100+) and Japanese (JLPT 8,500+). During study, it can optionally call the local Codex CLI / Claude Code to generate example sentences or polish definitions; it also runs offline without affecting the study path.

## Base Directory Structure

Create or maintain files in this structure. Do not create parallel directories for the same file type unless the project already has a stronger convention.

- `CLAUDE.md`: shared workflow for repository baseline, directory structure, after-change requirements, plan/review lifecycle, review expiry, file-size guardrail, project-specific triggers, project conventions, and validation.
- `AGENTS.md`: entry and tool differences; it references and follows the shared rules in `CLAUDE.md`.
- `UI_COPY_LANGUAGE.md`: SSOT for user-facing CLI copy language and locale mode.
- `README.md`: user and maintainer instructions for installation, usage, validation, and structure.
- `src/`: first-party Go source code and built-in vocabulary data. Vocabulary JSON under `src/data/` is committed to git and read-only at runtime.
- `scripts/`: project scripts and automation helpers.
- `build/`: local build artifacts. `Makefile` uses `go build -o build/vocabmaster ./src`.
- `dist/`: reserved as an optional packaging output root; there are currently no active artifacts.
- `ref/changelogs/INDEX.md`: final changelog index. Final changelog files use `ref/changelogs/CHANGELOG_X_<topic>.md`; existing historical `CHANGELOG_X.md` files keep their current names.
- `ref/reviews/INDEX.md`: final review index. Final review files use `ref/reviews/REVIEW_X_<topic>.md`; existing historical `REVIEW_X.md` files keep their current names.
- `ref/plans/INDEX.md`: final plan index. New final plan files use `ref/plans/PLAN_X_<topic>.md`; existing historical slug-date plans keep their current names.
- `ref/conventions/INDEX.md`: index of promoted project conventions. Convention bodies use `ref/conventions/CONVENTION_X_<topic>.md`.
- `ref/conventions/tally.md`: entry point for counting repeated feedback / repeated agent mistakes.
- `.ref/`: must be listed in `.gitignore`; holds only non-final plan/review working copies, not final records.

## UI/CLI Copy Language

Write active project documentation and maintainer/agent-facing instructions in English by default, including changelogs, plans, reviews, and conventions. Exceptions are `UI_COPY_LANGUAGE.md`, user-facing CLI copy governed by that file, locale examples, quoted/source text, and explicit non-English trigger anchors or examples.

Before adding or changing user-facing CLI copy, read `UI_COPY_LANGUAGE.md` and follow its active mode. If the requested copy language or supported locales differ from that file, update `UI_COPY_LANGUAGE.md` first, then make the CLI copy change.

## Required After Changes

Before starting, run `ls ref/conventions ref/changelogs ref/plans ref/reviews 2>/dev/null || true` from the repository root and review relevant records. Missing directories are setup work, not an error.

After changes, apply these minimum rules before any project-specific triggers:

1. If you change user-visible behavior, user-facing CLI copy, file structure, launch method, ports, dependencies, or validation steps, update the matching `README.md` section and follow `UI_COPY_LANGUAGE.md`. If the language requirements differ, update that file first. Pure bug fixes and internal refactors do not require README changes.
2. For every meaningful feature / behavior / command / dependency / structure change, write `ref/changelogs/CHANGELOG_X_<topic>.md` and update `ref/changelogs/INDEX.md`. For debug / performance / security / review-driven fixes, write `ref/reviews/REVIEW_X_<topic>.md` and update `ref/reviews/INDEX.md`. Before creating any final `ref/` record, run `ls <target-dir>/`, choose `X` as the next integer after the current maximum same-type number, and do not guess. Use a short stable kebab-case `<topic>` that is not vague like `update`, `fix`, or `misc`. INDEX summaries must be <= 80 characters or one short English sentence.
3. Keep non-final plan/review files in the current environment's workspace; use `<repo>/.ref/` when there is no stronger contract. At final closeout, archive the final plan to `ref/plans/PLAN_X_<topic>.md`, archive the final review to `ref/reviews/REVIEW_X_<topic>.md`, update the corresponding INDEX, and clean up the workspace copy.
4. Keep the advisory plan archive pre-commit hook installed with `bash scripts/plan-archive-reminder-pre-commit.sh --install` after setup or whenever `.git/hooks/pre-commit` is reset. The hook reminds about non-final `.ref/plans/` files and must not block commits.
5. Record repeated user feedback or repeated agent mistakes in `ref/conventions/tally.md` first. After `count >= 3`, run this repository's review flow, promote the rule to `ref/conventions/CONVENTION_X_<topic>.md`, and update `ref/conventions/INDEX.md`.

## Project-Specific Triggers

- Changes to vocabulary JSON under `src/data/`: follow the data update flow. Data-only refreshes do not enter the main changelog line; add review records depending on quality risk.
- Changes to CLI commands, installation method, build output path, or LLM provider order: update the corresponding `README.md` section and write a changelog.

Data update flow for committed vocabulary JSON:

1. Identify the source input: ECDICT SQLite for English, JLPT JSON for Japanese, or a documented replacement source. Record source path/version in the review record when data quality or counts change materially.
2. Regenerate with `python3 src/tools/extract_words.py` when using the default local inputs (`/tmp/ecdict-sqlite/stardict.db` and `/tmp/jlpt.json`), or document the equivalent replacement process before editing `src/data/english.json` or `src/data/japanese.json` manually.
3. Validate JSON shape before committing: top-level `version`, `language`, and `words`; each word keeps stable `id`, `language`, `text`, `chinese_def`, `difficulty`, `examples`, and optional pronunciation / part-of-speech / tags fields consistent with README's import format.
4. Check count and quality deltas against README expectations: English remains ECDICT-derived, Japanese remains JLPT-derived, difficulty bands stay meaningful, and obvious duplicates / malformed IDs / empty Chinese definitions are rejected.
5. Run `make test` after data changes; run `make build` as well when embedded data size, build output, or runtime loading behavior changes.
6. Data-only refreshes do not require a main changelog unless they change user-visible commands, schema, install/build behavior, or documented source policy. Write `ref/reviews/REVIEW_X_<topic>.md` and update `ref/reviews/INDEX.md` when the refresh depends on new source quality assumptions, count shifts, filtering logic, or manual curation risk.

## Project-Specific Conventions

> Dynamic upgrades go through `ref/conventions/CONVENTION_X_<topic>.md`; this section keeps only project invariants that must be visible at entry.

- Keep `go.mod` / `go.sum` at the repository root. First-party Go packages live under `src/`, and import paths use `github.com/vocabmaster/vocabmaster/src/...`.
- LLM enhancement is optional local functionality. The provider order is fixed as `codex -> claude -> fail`; when neither Codex CLI nor Claude Code is available or calls fail, the study path must continue using the built-in base data.

## Review Expiry And Minimum Re-Review Scope

Use this section to determine the minimum scope for the next review. `ref/reviews/` contains expiring coverage records, not permanent exemptions.

Minimum scope for the next review:

```text
unreviewed files ∪ expired reviewed files ∪ scope_unknown files
```

Since the most recent REVIEW baseline that covered a file, coverage expires when any of the following is true:

- Net changes >= `min(200 lines, 30% of current LOC)`.
- Distinct commit count >= 3.
- At least 90 days have passed and the file has changed at least once.
- REVIEW frontmatter marks `expired: true`.

When preparing a review, run `bash scripts/file-level-review-expiry.sh` at the repository root. If the script is missing, use `git log` to determine status manually according to the rules above.

## File-Size Guardrail (500 Lines)

If any source file exceeds 500 LOC, attempt to split it before submitting. Generated code, lockfiles, snapshots, migrations, and fixtures (including vocabulary JSON under `src/data/`) are excluded.

Splitting priority:

1. Extract module-level pure functions / types / constants.
2. Organize into same-directory submodules while preserving import paths.
3. Use facade + shared context splitting only after plan/review.

When a file truly cannot be split, record the file and the concrete reason in the relevant changelog's "do not split" protection list.

## Validation Flow

```bash
make build         # go build -> build/vocabmaster
make test          # go test ./... + installer shell tests
make install       # install to GOPATH/bin; must be tested when changing installation flow
make clean         # rm -rf build
```

After changing user-visible CLI behavior, installation flow, build path, or LLM provider order, run at least `make test`. For installation-related changes, also test `make install` or an equivalent isolated installation command.

## Deployment / Packaging

There is currently no separate deployment flow. Local installation and pre-release checks are based on `make install`, `make build`, and `make test`.
