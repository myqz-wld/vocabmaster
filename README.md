# VocabMaster

命令行背单词工具，支持英文和日文，基于 SM-2 间隔重复算法。

## 特性

- **20,000+ 内置词库** — 英文 12,100+（ECDICT）、日文 8,500+（JLPT N5-N1）
- **SM-2 间隔重复** — 基于遗忘曲线自动调度复习，科学记忆
- **一键学习** — `study` 命令自动平衡复习与新词，无需手动管理
- **中文释义 + 发音标注** — 英文含 IPA 音标，日文含假名读音
- **LLM 实时增强** — 学习和复习时自动调用本地 Claude Code 生成例句、润色释义
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

### 手动构建

```bash
git clone https://github.com/myqz-wld/vocabmaster.git
cd vocabmaster
make build
# 可执行文件在 ./vocabmaster
```

## 快速开始

```bash
# 一键学习（推荐，自动平衡复习和新词）
vocabmaster study

# 只学日文
vocabmaster study --lang ja

# 只学初级英文
vocabmaster study --lang en --level 1
```

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
vocabmaster learn --lang en --level 1 --count 3

# 复习到期单词
vocabmaster review

# 复习全部到期单词（不限数量）
vocabmaster review --count 0

# 查看统计
vocabmaster stats

# 搜索单词
vocabmaster search 环境

# 查看某个词的详情
vocabmaster info en_environment

# 浏览日文中级词库
vocabmaster list --lang ja --level 2

# 导入自定义词库
vocabmaster import my_words.json
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

学习新词或复习时，如果本地安装了 [Claude Code](https://claude.ai/claude-code)，会自动调用进行：

- 润色中文释义
- 生成自然例句（目标语言 + 中文翻译）
- 校验发音标注

结果缓存在本地数据库，每个词只调用一次。如果 Claude CLI 不可用，直接使用内置基础数据，不影响正常使用。

也可以通过 `generate` 命令批量预处理：

```bash
vocabmaster generate --lang en --count 100
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
