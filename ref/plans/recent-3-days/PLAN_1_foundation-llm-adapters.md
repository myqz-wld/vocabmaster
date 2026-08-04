---
plan_id: 1
status: completed
completed_at: 2026-08-04
---

# Foundation Alignment and Configurable LLM Adapters

## Goal

Align VocabMaster with the current `project-engineering-foundation` templates and add configurable Claude, Codex, and Grok adapters for LLM enrichment, including adapter-specific model and thinking-effort selection.

## Invariants

- Preserve the existing Go module path, `src/` source layout, `build/` output root, install metadata flow, CLI copy language, and offline study fallback.
- Keep historical reference record content intact while routing terminal records into the current time-bucket structure.
- Keep enrichment optional: unavailable or failed adapters must never block the normal study path.
- Keep existing cached enrichment behavior; `generate --force` remains the explicit way to regenerate cached entries.
- Keep all non-generated source files at or below 500 LOC by splitting the current LLM implementation before extending it.

## Confirmed Scope and Design Decisions

- Use existing-repository `minimal repair`, merging current foundation template rules without replacing project-specific invariants.
- Add all four current time buckets and indexes under `ref/changelogs/`, `ref/reviews/`, and `ref/plans/`; legacy records with missing or old dates route to `history/`.
- Replace the old plan-only archive reminder with the current `.ref` archive reminder and update the review-expiry helper.
- Add persistent CLI flags: `--llm-adapter`, `--llm-model`, and `--llm-thinking`.
- Default adapter mode is `auto`, ordered `codex -> claude -> grok`; explicit adapter mode tries only the selected adapter.
- Model and thinking values are optional and use the selected CLI's defaults when omitted. They require an explicit adapter because model identifiers and supported effort levels are provider-specific.
- Pass model/thinking through native CLI arguments: Claude `--model` / `--effort`, Codex `--model` / `model_reasoning_effort`, and Grok `--model` / `--reasoning-effort`.

## Checklist

- [x] Align `CLAUDE.md`, `AGENTS.md`, `UI_COPY_LANGUAGE.md`, `.gitignore` checks, scripts, and `ref/` indexes with current foundation templates.
- [x] Rebucket historical changelogs, reviews, and plans without rewriting their historical content.
- [x] Split `src/internal/llm/llm.go` into focused client, adapter, prompt, and response modules.
- [x] Implement configurable Claude, Codex, and Grok adapter execution.
- [x] Wire CLI options through batch generation and interactive study/review/learn sessions.
- [x] Add adapter/options/argument/fallback tests and keep test files below 500 LOC.
- [x] Update README and project workflow invariants.
- [x] Run formatting, unit tests, full project tests, build, CLI help smoke tests, and file-size checks.
- [x] Archive this plan and the feature/foundation changelog in current time buckets; refresh all indexes.
- [x] Install and exercise the advisory `.ref` archive pre-commit hook.

## Risks and Validation Requirements

- CLI flags can drift across third-party releases; unit tests must lock the locally verified argument shape and README must state that model/effort support is adapter/model dependent.
- Grok output must stay plain final text so the existing robust JSON cleanup pipeline receives the model response rather than a headless metadata envelope.
- Explicit adapter selection must affect both availability checks and execution; tests must prove no fallback to unselected adapters.
- `auto` plus model/thinking must fail before processing any words to avoid repeated per-word CLI failures.
- Moving legacy records can break relative links; scan Markdown links after rebucketing.

## Validation

- `git diff --check`
- `~/.local/bin/mise exec -- go test ./...`
- `~/.local/bin/mise exec -- go test -race ./src/internal/llm`
- `make test`
- `make build`
- Isolated install, installed `vm --version`, and uninstall cycle
- Shell syntax, CLI help, invalid configuration, archive reminder, review-expiry, and file-size checks
- Local Claude, Codex, and Grok argument-contract probes using `--help` only; no model request was made

## Final Status

Completed on 2026-08-04. The foundation layout, reference routing, helper scripts, concise user README, and configurable enrichment adapters are implemented and validated. Related record: `CHANGELOG_13_foundation-llm-adapters.md`.
