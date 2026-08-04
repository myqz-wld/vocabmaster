# VocabMaster

A terminal vocabulary trainer for English and Japanese, powered by SM-2 spaced repetition.

- 20,000+ built-in words from ECDICT and JLPT
- Chinese definitions, pronunciations, and example sentences
- Automatic review scheduling and one-command study sessions
- Optional local LLM enrichment through Codex, Claude, or Grok
- Offline-first: study continues normally when no LLM is available

## Install

Requires Go 1.24 or newer.

```bash
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make install
```

The installer adds `vocabmaster` and the short command `vm` to `$(go env GOPATH)/bin`. Reopen your terminal or run the printed `source` command once installation completes.

```bash
# Optional custom install directory
make install BINDIR=/path/to/bin

# Do not edit shell configuration automatically
make install UPDATE_SHELL_RC=0
```

## Quick Start

```bash
vm study                    # Review due words, then learn new ones
vm study --lang ja          # Japanese only
vm learn --lang en --count 5
vm review                   # Review due words
vm stats                    # Show learning statistics
vm search 环境
```

`study` automatically adjusts new-word volume to the current review load. Use `--new-words` to override it.

## Commands

| Command | Purpose |
|---|---|
| `study` | Automatically balance reviews and new words |
| `learn` | Learn new words only |
| `review` | Review due words only |
| `stats` | Show learning statistics |
| `list` / `search` / `info` | Browse the vocabulary database |
| `import` | Import a JSON vocabulary file |
| `generate` | Batch-generate LLM enrichment |
| `config` | View or save local LLM defaults |
| `reset` | Reset learning progress |
| `--version` | Show installed build information |
| `--check-installed` | Check whether the installed CLI matches this checkout |

Run `vm <command> --help` for all options.

## LLM Enrichment

VocabMaster can use locally installed Codex, Claude Code, or Grok CLIs to polish definitions, validate pronunciations, and generate example sentences.

Default `auto` mode tries `codex -> claude -> grok`. Selecting an adapter explicitly disables fallback to the others.

```bash
# Use Claude with its defaults
vm study --llm-adapter claude

# Choose a Codex model and reasoning effort
vm generate --lang en --count 100 \
  --llm-adapter codex \
  --llm-model gpt-5.6 \
  --llm-thinking high

# Rebuild cached enrichment with Grok
vm generate --lang ja --count 100 --force \
  --llm-adapter grok \
  --llm-model grok-4.5 \
  --llm-thinking high
```

| Adapter | Accepted `--llm-thinking` values |
|---|---|
| `codex` | `low`, `medium`, `high`, `xhigh`, `max` |
| `claude` | `low`, `medium`, `high`, `xhigh`, `max` |
| `grok` | `low`, `medium`, `high`, `xhigh` |

`--llm-model` and `--llm-thinking` require an explicit `--llm-adapter`. The selected model must support the requested effort.

Save long-lived defaults on this machine with:

```bash
vm config set-llm --adapter codex --model gpt-5.6-luna --thinking high
vm config show
```

The settings are stored in `~/.vocabmaster/config.json`. Command-line `--llm-*` options temporarily override them.

Enrichment is cached per word. The options affect cache misses; use `generate --force` to regenerate existing entries with a different adapter, model, or effort.

## Vocabulary

| Language | Source | Levels |
|---|---|---|
| English | ECDICT | Beginner, intermediate, advanced |
| Japanese | JLPT N5-N1 | Beginner, intermediate, advanced |

Import custom words with `vm import <file.json>`. Each file contains `version`, `language`, and a `words` array; each word uses fields such as `id`, `text`, `chinese_def`, `difficulty`, `pronunciation`, `part_of_speech`, `examples`, and `tags`.

## Data

Learning progress and cached enrichment are stored in `~/.vocabmaster/`. Use `--data-dir` to choose another location.

## Acknowledgements

- [ECDICT](https://github.com/skywind3000/ECDICT)
- [JLPT Vocabulary](https://github.com/Bluskyo/JLPT_Vocabulary)
- SuperMemo SM-2 algorithm
