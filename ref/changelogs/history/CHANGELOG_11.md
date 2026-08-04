# CHANGELOG_11

## Summary

Tightened foundation prompt wording and the durable LLM enrichment prompt contract.

## Changes

- Aligned `AGENTS.md` with the current foundation template so it only points to `CLAUDE.md` and entry-specific differences.
- Tightened `CLAUDE.md` foundation wording toward the current template while preserving VocabMaster-specific runtime, layout, zh-CN CLI copy, review expiry, file-size, validation, and LLM provider-order rules.
- Strengthened the source-held LLM enrichment prompt to preserve the word, part of speech, and target-language meaning, and to require target-language example sentences with natural Chinese translations.
- Kept the enrichment prompt's Chinese terminology consistent by using `单词` throughout the word data block and request line.

## Validation

- `git diff --check`
- `go test ./src/internal/llm`
