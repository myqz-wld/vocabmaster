# VocabMaster

A command-line vocabulary memorization tool for English and Japanese, based on the SM-2 spaced repetition algorithm.

## Features

- **20,000+ built-in vocabulary entries** — 12,100+ English entries (ECDICT) and 8,500+ Japanese entries (JLPT N5-N1)
- **SM-2 spaced repetition** — automatically schedules reviews based on the forgetting curve for efficient memorization
- **One-command study** — the `study` command automatically balances review items and new words, with no manual management required
- **Chinese definitions + pronunciation notes** — English entries include IPA pronunciations; Japanese entries include kana readings
- **Real-time LLM enhancement** — during study and review, automatically calls the local Codex CLI / Claude Code to generate example sentences and polish definitions
- **Three difficulty levels** — beginner / intermediate / advanced, classified by Oxford 3000, Collins stars, and JLPT level
- **Custom imports** — supports importing external JSON vocabulary files

## Installation

### Install From Source

```bash
# Requires Go 1.24+
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make install
```

By default, `make install` installs to `$(go env GOPATH)/bin`, creates the short command `vm`, and writes the install directory and `vm` alias to the current shell configuration file:

```bash
vm study
```

After the first installation, run `source <config-file>` as printed by the installer, or reopen the terminal for the new environment to take effect. The zsh default is:

```bash
source ~/.zshrc
```

Automatic shell configuration updates support zsh, bash, and POSIX profile syntax. For other shells, skip the automatic update and configure manually.

To install to another directory:

```bash
make install BINDIR=/path/to/bin
```

To prevent the installer from modifying shell configuration:

```bash
make install UPDATE_SHELL_RC=0
```

### Manual Build

```bash
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make build
# The executable is at ./build/vocabmaster
```

## Project Structure

- `src/`: Go source code and built-in vocabulary data
- `build/`: local build artifacts, not committed to git
- `ref/`: changelog, review, plan, and project convention records
- `scripts/`: repository maintenance scripts (review expiry checks)

## Quick Start

```bash
# One-command study (recommended; automatically balances reviews and new words)
vm study

# Study Japanese only
vm study --lang ja

# Study beginner English only
vm study --lang en --level 1
```

After installation, both `vm` and `vocabmaster` commands are available. The examples below use `vm`.

## Commands

| Command | Description |
|------|------|
| `study` | One-command study (automatically balances reviews and new words) |
| `learn` | Study new words only |
| `review` | Review due words only |
| `stats` | View learning statistics |
| `list` | Browse the vocabulary database |
| `search` | Search words/definitions (AI-enhanced data is shown first) |
| `info` | View word details and learning progress |
| `import` | Import an external JSON vocabulary file |
| `generate` | Batch LLM preprocessing for the vocabulary database |
| `reset` | Reset learning progress |

## Usage Examples

```bash
# Learn 3 new beginner English words
vm learn --lang en --level 1 --count 3

# Review due words
vm review

# Review all due words (no limit)
vm review --count 0

# View statistics
vm stats

# Search for a word
vm search 环境

# View details for a specific word
vm info en_environment

# Browse intermediate Japanese vocabulary
vm list --lang ja --level 2

# Import a custom vocabulary file
vm import my_words.json
```

## Smart Scheduling In The `study` Command

`study` automatically decides what to do based on the current study load:

| Due Review Count | Behavior |
|-----------|------|
| > 20 words | Focus on review; do not learn new words |
| 11-20 words | Review first, then learn 5 new words |
| <= 10 words | Review first, then learn 10 new words |

Each review session includes up to 30 due words. Use `--new-words` to override the default behavior.

## Vocabulary Levels

### English (Source: ECDICT)

| Level | Criteria | Count |
|------|------|------|
| Beginner | Oxford 3000 / Collins 4-5 stars / high-frequency words | ~3,000 |
| Intermediate | Collins 3 stars / medium-frequency words / CET-4/6 | ~4,100 |
| Advanced | Collins 1-2 stars / GRE / TOEFL / IELTS | ~5,000 |

### Japanese (Source: JLPT)

| Level | Criteria | Count |
|------|------|------|
| Beginner | N5 + N4 | ~1,350 |
| Intermediate | N3 | ~1,800 |
| Advanced | N2 + N1 | ~5,300 |

## LLM Enhancement

When learning new words or reviewing, if the Codex CLI or [Claude Code](https://claude.ai/claude-code) is installed locally, VocabMaster automatically calls them in `codex -> claude` order to:

- Polish Chinese definitions
- Generate natural example sentences (target language + Chinese translation)
- Validate pronunciation notes

Results are cached in the local database, and each word is processed only once. If neither Codex CLI nor Claude CLI is available, or calls fail, VocabMaster uses the built-in base data directly and normal usage is unaffected.

You can also preprocess in batches with the `generate` command:

```bash
vm generate --lang en --count 100
```

## Custom Vocabulary Format

Imported JSON files must use the following format:

```json
{
  "version": "1.0",
  "language": "en",
  "words": [
    {
      "id": "en_example",
      "language": "en",
      "text": "example",
      "pronunciation": "/ɪɡˈzæm.pəl/",
      "chinese_def": "例子；示例",
      "difficulty": 2,
      "part_of_speech": "noun",
      "examples": [
        {
          "sentence": "This is a good example.",
          "translation": "这是一个好例子。"
        }
      ],
      "tags": ["education"]
    }
  ]
}
```

## Data Storage

- Learning progress: `~/.vocabmaster/vocabmaster.db` (SQLite)
- Use `--data-dir` to specify another directory

## Acknowledgements

- [ECDICT](https://github.com/skywind3000/ECDICT) — English-Chinese dictionary data
- [JLPT_Vocabulary](https://github.com/Bluskyo/JLPT_Vocabulary) — JLPT Japanese vocabulary data
- SM-2 algorithm — SuperMemo spaced repetition algorithm
