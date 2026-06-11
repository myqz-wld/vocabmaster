# VocabMaster

命令行背单词工具，支持英文和日文，基于 SM-2 间隔重复算法。

## 特性

- **20,000+ 内置词库** — 英文 12,100+（ECDICT）、日文 8,500+（JLPT N5-N1）
- **SM-2 间隔重复** — 基于遗忘曲线自动调度复习，科学记忆
- **一键学习** — `study` 命令自动平衡复习与新词，无需手动管理
- **中文释义 + 发音标注** — 英文含 IPA 音标，日文含假名读音
- **LLM 实时增强** — 学习和复习时自动调用本地 Codex CLI / Claude Code 生成例句、润色释义
- **三级难度** — 初级 / 中级 / 高级，按 Oxford 3000、Collins 星级、JLPT 等级分类
- **自定义导入** — 支持导入外部 JSON 词库

## 安装

### 从源码安装

```bash
# 需要 Go 1.24+
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make install
```

`make install` 默认安装到 `$(go env GOPATH)/bin`，创建 `vm` 短命令，并把安装目录和 `vm` alias 写入当前 shell 的配置文件：

```bash
vm study
```

首次安装后按安装脚本输出运行 `source <配置文件>`，或重开终端让新环境生效。zsh 默认是：

```bash
source ~/.zshrc
```

自动写入 shell 配置支持 zsh、bash 和 POSIX profile 语法；其它 shell 可跳过自动写入后手动配置。

如果需要安装到其他目录：

```bash
make install BINDIR=/path/to/bin
```

如果不希望安装脚本修改 shell 配置：

```bash
make install UPDATE_SHELL_RC=0
```

### 手动构建

```bash
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make build
# 可执行文件在 ./build/vocabmaster
```

## 项目结构

- `src/`：Go 源码、内置词库数据和维护脚本
- `build/`：本地构建产物，不入 git
- `ref/`：changelog、review、plan 和项目约定记录
- `scripts/`：仓库维护脚本（review 过期检查）

## 快速开始

```bash
# 一键学习（推荐，自动平衡复习和新词）
vm study

# 只学日文
vm study --lang ja

# 只学初级英文
vm study --lang en --level 1
```

安装后 `vm` 和 `vocabmaster` 两个命令都可用；下方示例使用 `vm`。

## 命令

| 命令 | 说明 |
|------|------|
| `study` | 一键学习（自动平衡复习和新词） |
| `learn` | 仅学习新词 |
| `review` | 仅复习到期单词 |
| `stats` | 查看学习统计 |
| `list` | 浏览词库 |
| `search` | 搜索单词/释义（优先展示 AI 增强数据） |
| `info` | 查看单词详情和学习进度 |
| `import` | 导入外部 JSON 词库 |
| `generate` | 批量 LLM 预处理词库 |
| `reset` | 重置学习进度 |

## 使用示例

```bash
# 学习 3 个新的初级英文单词
vm learn --lang en --level 1 --count 3

# 复习到期单词
vm review

# 复习全部到期单词（不限数量）
vm review --count 0

# 查看统计
vm stats

# 搜索单词
vm search 环境

# 查看某个词的详情
vm info en_environment

# 浏览日文中级词库
vm list --lang ja --level 2

# 导入自定义词库
vm import my_words.json
```

## study 命令的智能调度

`study` 会根据当前学习负载自动决定：

| 待复习词数 | 行为 |
|-----------|------|
| > 20 词 | 专注复习，不学新词 |
| 10-20 词 | 先复习，再学 5 个新词 |
| < 10 词 | 先复习，再学 10 个新词 |

每次复习最多 30 个到期词。可通过 `--new-words` 覆盖默认行为。

## 词库分级

### 英文（来源：ECDICT）

| 级别 | 标准 | 数量 |
|------|------|------|
| 初级 | Oxford 3000 / Collins 4-5 星 / 高频词 | ~3,000 |
| 中级 | Collins 3 星 / 中频 / CET-4/6 | ~4,100 |
| 高级 | Collins 1-2 星 / GRE / 托福 / 雅思 | ~5,000 |

### 日文（来源：JLPT）

| 级别 | 标准 | 数量 |
|------|------|------|
| 初级 | N5 + N4 | ~1,350 |
| 中级 | N3 | ~1,800 |
| 高级 | N2 + N1 | ~5,300 |

## LLM 增强

学习新词或复习时，如果本地安装了 Codex CLI 或 [Claude Code](https://claude.ai/claude-code)，会按 `codex -> claude` 顺序自动调用进行：

- 润色中文释义
- 生成自然例句（目标语言 + 中文翻译）
- 校验发音标注

结果缓存在本地数据库，每个词只调用一次。如果 Codex CLI 和 Claude CLI 都不可用或调用失败，直接使用内置基础数据，不影响正常使用。

也可以通过 `generate` 命令批量预处理：

```bash
vm generate --lang en --count 100
```

## 自定义词库格式

导入的 JSON 文件需符合以下格式：

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

## 数据存储

- 学习进度：`~/.vocabmaster/vocabmaster.db`（SQLite）
- 可通过 `--data-dir` 指定其他目录

## 致谢

- [ECDICT](https://github.com/skywind3000/ECDICT) — 英汉词典数据
- [JLPT_Vocabulary](https://github.com/Bluskyo/JLPT_Vocabulary) — JLPT 日语词汇数据
- SM-2 算法 — SuperMemo 间隔重复算法
