---
plan_id: src-layout-alignment-20260608
status: completed
created_at: 2026-06-08T06:42:30Z
base_branch: main
base_commit: ef298fc
worktree_path: /Users/wanglidong/Repository/vocabmaster
final_commit: 3365efe157927bbfb8cfce912cc54324a08494a0
completed_at: 2026-06-08T14:57:52+08:00
---

# Plan: align VocabMaster repository layout with project-engineering-foundation

## Goal

Make the repository match the skill layout: first-party source under `src/`, build output under `build/`, root metadata at the repository root, and durable AI-coding entry points in `CLAUDE.md` / `AGENTS.md`.

## Invariants

- `go.mod` and `go.sum` stay at the repository root.
- `build/` remains the only active build output root; `dist/` stays ignored.
- Historical records in `ref/plans/` and old changelogs keep their original paths as historical facts.
- `make build`, `make test`, `make install`, and `vm --help` must work after the move.

## Design Decisions

- Move `main.go`, `cmd/`, `internal/`, `data/`, and `tools/` under `src/`.
- Build from `./src` while keeping the module path `github.com/vocabmaster/vocabmaster`.
- Update imports to `github.com/vocabmaster/vocabmaster/src/...`.
- Create `AGENTS.md` as the runtime-specific companion to `CLAUDE.md`; keep shared repository rules in `CLAUDE.md`.

## Checklist

- [x] Move first-party source/data/tools under `src/`.
- [x] Update Go imports and Makefile paths.
- [x] Update `CLAUDE.md`, create `AGENTS.md`, and update `README.md`.
- [x] Add changelog record for the structure change.
- [x] Run validation: `go test ./...`, `make test`, `make build`, `./build/vocabmaster --help`, isolated install test, real `make install`.

## Current Progress

Source/data/tools moved under `src/`; imports, Makefile, README, CLAUDE.md, AGENTS.md, and CHANGELOG_4 updated. Validation and simple-review passed.

## Known Risks

- Go `internal/` import visibility changes when moved to `src/internal`; packages under `src/` remain inside the allowed parent.
- Historical refs contain old paths and must not be rewritten as if they were current facts.

## Next-Session First Action

From `/Users/wanglidong/Repository/vocabmaster`, inspect this plan and run `git status --short`, then continue at the first unchecked checklist item.
