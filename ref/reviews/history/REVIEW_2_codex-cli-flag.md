# REVIEW_2: Codex CLI flag compatibility

## Trigger

User reported that LLM enrichment appeared to call Claude first despite the documented `codex -> claude -> fail` provider order.

## Method

- Checked source provider order in `src/internal/llm/llm.go`.
- Ran provider-order unit tests for Codex-first and Claude fallback behavior.
- Compared installed `vm` / `vocabmaster` build metadata with the current checkout.
- Reproduced installed-binary provider order using temporary fake `codex` and `claude` executables.
- Smoke-tested the real local Codex CLI argument set.

## Triage

| Finding | Decision | Handling |
|---|---|---|
| Source and installed binary still use `codex -> claude -> fail` order | ✅ confirmed | No provider-order change required |
| Installed CLI matches current checkout commit `f17f19f` | ✅ confirmed | Ruled out stale install as the cause |
| Current `codex-cli 0.142.0` rejects `--ask-for-approval never` | ✅ P2 | Replaced the obsolete flag with `-c approval_policy="never"` |
| Fast Codex startup failure makes Claude appear to run first | ✅ explained | Codex failed before visible LLM work, then normal fallback invoked Claude |

## Fixes

- Updated `runCodex` to build arguments through `codexExecArgs`.
- Replaced the obsolete `--ask-for-approval never` flag with `-c approval_policy="never"`.
- Kept `--sandbox read-only`, `--skip-git-repo-check`, `--ephemeral`, `--ignore-rules`, `--color never`, and `--output-last-message`.
- Added regression coverage asserting the current Codex argument list and preventing the obsolete flag from returning.

## Verification

- `go test ./src/internal/llm -run 'TestEnrichWordUsesCodexFirst|TestEnrichWordFallsBackToClaude|TestIsAvailableChecksAnyProvider' -v`
- Temporary fake-provider installed-binary check:
  - Codex success: provider order `codex`
  - Codex failure: provider order `codex -> claude`
- Real Codex smoke:
  - `codex exec -c 'approval_policy="never"' --sandbox read-only --skip-git-repo-check --ephemeral --ignore-rules --color never --output-last-message <tmp> 'Return exactly: ok'`
  - Exit `0`; output showed `approval: never`, `sandbox: read-only`, and final message `ok`.
- `go test ./src/internal/llm -v`

## Residual Info

- Claude CLI on this machine returned an API-disabled result during manual probing; that is separate from provider order and does not affect the Codex flag fix.
