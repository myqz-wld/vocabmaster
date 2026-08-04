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

A command-line vocabulary memorization tool based on the SM-2 spaced repetition algorithm. It supports English (ECDICT 12,100+) and Japanese (JLPT 8,500+). During study, it can optionally call local Codex, Claude, or Grok CLI adapters to generate example sentences or polish definitions; it also runs offline without affecting the study path.

## Base Directory Structure

Create or maintain files in this structure. Do not create parallel directories for the same file type unless the project already has a stronger convention.

- `CLAUDE.md`: shared workflow for repository baseline, directory structure, after-change requirements, plan/review lifecycle, review expiry, file-size guardrail, project-specific triggers, project conventions, and validation.
- `AGENTS.md`: entry and tool differences; it references and follows the shared rules in `CLAUDE.md`.
- `UI_COPY_LANGUAGE.md`: SSOT for user-facing CLI copy language and locale mode.
- `README.md`: user and maintainer instructions for installation, usage, validation, and structure.
- `src/`: first-party Go source code and built-in vocabulary data. Vocabulary JSON under `src/data/` is committed to git and read-only at runtime.
- `scripts/`: project scripts and automation helpers, including installation, packaging, data extraction, build metadata, review expiry, and `.ref` archive reminders.
- `build/`: local build artifacts. `Makefile` uses `go build -o build/vocabmaster ./src`.
- `dist/`: local distributable archives and checksums produced by `make pack`.
- `ref/changelogs/INDEX.md`: final changelog routing index. New final changelogs use `ref/changelogs/<bucket>/CHANGELOG_X_<topic>.md`; historical filenames remain unchanged when rebucketed.
- `ref/reviews/INDEX.md`: final review routing index. New final reviews use `ref/reviews/<bucket>/REVIEW_X_<topic>.md`; historical filenames remain unchanged when rebucketed.
- `ref/plans/INDEX.md`: final plan routing index. New final plans use `ref/plans/<bucket>/PLAN_X_<topic>.md`; historical filenames remain unchanged when rebucketed.
- `ref/*/{recent-3-days,recent-week,recent-month,history}/INDEX.md`: mutually exclusive time-bucket indexes for final changelogs, reviews, and plans.
- `ref/conventions/INDEX.md`: index of promoted project conventions. Convention bodies use `ref/conventions/CONVENTION_X_<topic>.md`.
- `ref/conventions/tally.md`: entry point for counting repeated feedback / repeated agent mistakes.
- `.ref/`: must be listed in `.gitignore`; holds non-final plans, reviews, raw outputs, spike drafts, scratch notes, and other unarchived LLM-facing material, never final records.

## UI/CLI Copy Language

Write active project documentation and maintainer/agent-facing instructions in English by default, including changelogs, plans, reviews, and conventions. Exceptions are `UI_COPY_LANGUAGE.md`, user-facing CLI copy governed by that file, locale examples, quoted/source text, and explicit non-English trigger anchors or examples.

Before adding or changing user-facing CLI copy, read `UI_COPY_LANGUAGE.md` and follow its active mode. If the requested copy language or supported locales differ from that file, update `UI_COPY_LANGUAGE.md` first, then make the CLI copy change.

## Required After Changes

Before starting, run `find ref/changelogs ref/plans ref/reviews -maxdepth 2 -type f -name '*.md' 2>/dev/null || true` from the repository root and review relevant records. Also inspect `ref/conventions/` when project-specific conventions may apply. Missing directories are setup work, not an error. Before creating or moving any final typed `ref/` record, read the relevant root and bucket `INDEX.md` files, scan every same-type bucket, choose `X` as the maximum existing same-type number plus 1, and do not guess. Use a short stable kebab-case topic that is not vague like `update`, `fix`, or `misc`.

After changes, apply these minimum rules before any project-specific triggers:

1. If you change user-visible behavior, user-facing CLI copy, file structure, launch method, ports, dependencies, or validation steps, update the matching `README.md` section and follow `UI_COPY_LANGUAGE.md`. If the language requirements differ, update that file first. Pure bug fixes and internal refactors do not require README changes.
2. For every meaningful feature, behavior, command, dependency, or structure change, write `ref/changelogs/<bucket>/CHANGELOG_X_<topic>.md`, rebucket all changelogs by `changed_at`, and update the root routing index plus every affected bucket index. For debug, performance, security, or review-driven fixes, do the equivalent under `ref/reviews/` using `reviewed_at`. Keep index summaries to 80 characters or one short English sentence.
3. Keep non-final plans in the current environment's plan workspace; use `<repo>/.ref/plans/` when there is no stronger contract. Keep non-final review drafts and raw output in the current review workspace or `<repo>/.ref/reviews/`. At final closeout, archive plans and reviews into the correct time bucket, rebucket all same-type records, update the root and affected bucket indexes, and clean up workspace copies.
4. Store durable LLM-facing extra materials such as investigation notes, architecture notes, and reusable evidence somewhere under `ref/` and link them from the relevant final record. Keep temporary scratch and raw logs in `.ref/` or the current environment workspace.
5. Keep the advisory `.ref` archive hook installed with `bash scripts/ref-archive-reminder-pre-commit.sh --install` after setup or whenever `.git/hooks/pre-commit` is reset. It classifies unarchived `.ref/` files and must remain advisory with exit status 0.
6. Record repeated user feedback or repeated agent mistakes in `ref/conventions/tally.md` first. After `count >= 3`, run this repository's review flow, promote the rule to `ref/conventions/CONVENTION_X_<topic>.md`, and update `ref/conventions/INDEX.md`.

## Project-Specific Triggers

- Changes to vocabulary JSON under `src/data/`: follow the data update flow. Data-only refreshes do not enter the main changelog line; add review records depending on quality risk.
- Changes to CLI commands, installation method, build output path, or LLM adapter selection/order: update the corresponding `README.md` section and write a changelog.

Data update flow for committed vocabulary JSON:

1. Identify the source input: ECDICT SQLite for English, JLPT JSON for Japanese, or a documented replacement source. Record source path/version in the review record when data quality or counts change materially.
2. Regenerate with `python3 scripts/extract_words.py` when using the default local inputs (`/tmp/ecdict-sqlite/stardict.db` and `/tmp/jlpt.json`), or document the equivalent replacement process before editing `src/data/english.json` or `src/data/japanese.json` manually.
3. Validate JSON shape before committing: top-level `version`, `language`, and `words`; each word keeps stable `id`, `language`, `text`, `chinese_def`, `difficulty`, `examples`, and optional pronunciation / part-of-speech / tags fields consistent with README's import format.
4. Check count and quality deltas against README expectations: English remains ECDICT-derived, Japanese remains JLPT-derived, difficulty bands stay meaningful, and obvious duplicates / malformed IDs / empty Chinese definitions are rejected.
5. Run `make test` after data changes; run `make build` as well when embedded data size, build output, or runtime loading behavior changes.
6. Data-only refreshes do not require a main changelog unless they change user-visible commands, schema, install/build behavior, or documented source policy. Write `ref/reviews/<bucket>/REVIEW_X_<topic>.md` and update the root plus affected bucket indexes when the refresh depends on new source quality assumptions, count shifts, filtering logic, or manual curation risk.

## Project-Specific Conventions

> Dynamic upgrades go through `ref/conventions/CONVENTION_X_<topic>.md`; this section keeps only project invariants that must be visible at entry.

- Keep `go.mod` / `go.sum` at the repository root. First-party Go packages live under `src/`, and import paths use `github.com/vocabmaster/vocabmaster/src/...`.
- LLM enhancement is optional local functionality. Default `auto` mode tries `codex -> claude -> grok -> fail`; an explicit adapter selection tries only that adapter and may pass adapter-specific model and thinking settings. When no selected CLI is available or calls fail, the study path must continue using built-in base data.

## Review Expiry And Minimum Re-Review Scope

Use this section to determine the minimum scope for the next review. `ref/reviews/` contains expiring coverage records, not permanent exemptions.

The next review's minimum scope is:

```text
unreviewed files union expired reviewed files union scope_unknown files
```

`scope_unknown files` are files whose previous coverage cannot be trusted because the review lacks a parseable `review-scope`, lacks a usable `baseline_commit`, or cannot be mapped to the current path.

Since the latest REVIEW `baseline_commit` that covered a file, coverage expires when any of the following is true:

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

When a file truly cannot be split, record the path, concrete reason, and revisit trigger in the relevant changelog's "Do Not Split Protection" or review's "Residual Risk" section.

## Validation Flow

```bash
make build         # go build -> build/vocabmaster
make test          # go test ./... + installer shell tests
make pack          # build -> dist/<release>-<platform>.tar.gz + SHA-256
make install       # install to GOPATH/bin; must be tested when changing installation flow
make clean         # rm -rf build
```

After changing user-visible CLI behavior, installation flow, build path, or LLM provider order, run at least `make test`. For installation-related changes, also test `make install` or an equivalent isolated installation command.

## Deployment / Packaging

There is currently no separate deployment flow. VocabMaster is an installable CLI through `make install`; `make pack` creates a platform-specific archive and SHA-256 checksum under `dist/`. Local installation and pre-release checks are based on `make install`, `make pack`, `make build`, and `make test`.

Packaging or installation changes must generate and ship build metadata with the installed CLI. The metadata must include at least the app/package name, semantic version when available, full git commit, short git commit, branch when available, dirty flag when determinable, and build timestamp. The installed CLI must expose both a human-readable version/status entry and a machine-checkable freshness command or equivalent, such as `vocabmaster --version` / `vm --version` and `vocabmaster --check-installed` / `vm --check-installed`.

The freshness check compares installed metadata with the current source checkout commit, may compare local `origin/main`, must not fetch remotes, and must report missing installed metadata separately from a commit mismatch.
