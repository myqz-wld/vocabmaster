package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/model"
)

type EnrichResult struct {
	ChineseDef    string         `json:"chinese_def"`
	Pronunciation string         `json:"pronunciation"`
	Examples      []model.Example `json:"examples"`
	LinkedWordIDs []string       `json:"linked_word_ids"`
}

func EnrichWord(word *model.Word, allWordIDs []string) (*model.Word, error) {
	if word == nil {
		return nil, fmt.Errorf("word is nil")
	}

	prompt := buildEnrichPrompt(word, allWordIDs)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude", "-p", prompt, "--output-format", "json")
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("调用 Claude CLI 超时（60秒）")
		}
		return nil, fmt.Errorf("调用 Claude CLI 失败: %w", err)
	}

	var response struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("解析 Claude 响应失败: %w", err)
	}

	if strings.TrimSpace(response.Result) == "" {
		return nil, fmt.Errorf("AI 返回了空响应")
	}

	resultJSON := cleanJSONResponse(response.Result)
	resultJSON = fixControlCharsInStrings(resultJSON)
	resultJSON = fixTrailingCommas(resultJSON)

	var result EnrichResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		preview := resultJSON
		if len(preview) > 300 {
			preview = preview[:300] + "..."
		}
		return nil, fmt.Errorf("解析增强数据失败: %w\n清理后JSON: %s", err, preview)
	}

	enriched := *word
	if result.ChineseDef != "" {
		enriched.ChineseDef = result.ChineseDef
	}
	if result.Pronunciation != "" {
		enriched.Pronunciation = result.Pronunciation
	}
	if len(result.Examples) > 0 {
		enriched.Examples = result.Examples
	}
	if len(result.LinkedWordIDs) > 0 {
		enriched.LinkedWordIDs = result.LinkedWordIDs
	}
	return &enriched, nil
}

func buildEnrichPrompt(word *model.Word, allWordIDs []string) string {
	var langDesc string
	if word.Language == model.LangEnglish {
		langDesc = "英文"
	} else {
		langDesc = "日文"
	}

	otherLangIDs := filterOtherLangIDs(word, allWordIDs)

	return fmt.Sprintf(`你是一个专业的语言学家。请为以下%s单词提供增强数据。

单词信息：
- 单词: %s
- 语言: %s
- 当前中文释义: %s
- 当前发音: %s
- 词性: %s

可关联的其他语言词汇ID列表：%s

请严格按照以下JSON格式返回（不要包含其他文字）：
{
  "chinese_def": "准确、地道的中文释义",
  "pronunciation": "校验后的发音标注（英文用IPA音标，日文用假名）",
  "examples": [
    {"sentence": "自然的例句1", "translation": "中文翻译1"},
    {"sentence": "自然的例句2", "translation": "中文翻译2"}
  ],
  "linked_word_ids": ["如果有含义对应的其他语言词汇，填入其ID"]
}

要求：
1. 中文释义要准确且易懂
2. 例句要自然常用，不要太复杂
3. 关联词只选择含义确实对应的词
4. 如果找不到合适的关联词，linked_word_ids 返回空数组
5. 只返回纯JSON，不要用markdown代码块包裹，不要有多余文字
6. JSON必须严格合法：不要有尾部逗号，字符串内不要有换行`, langDesc, word.Text, word.Language, word.ChineseDef, word.Pronunciation, word.PartOfSpeech, strings.Join(otherLangIDs, ", "))
}

func cleanJSONResponse(raw string) string {
	s := strings.TrimSpace(raw)

	// 移除 markdown 代码块标记（支持代码块出现在任意位置）
	if idx := strings.Index(s, "```"); idx >= 0 {
		// 找到代码块开始后的第一个换行
		afterOpen := s[idx+3:]
		nlIdx := strings.Index(afterOpen, "\n")
		if nlIdx >= 0 {
			inner := afterOpen[nlIdx+1:]
			// 找到结束的 ```
			if closeIdx := strings.Index(inner, "```"); closeIdx >= 0 {
				s = strings.TrimSpace(inner[:closeIdx])
			} else {
				s = strings.TrimSpace(inner)
			}
		}
	}

	// 提取 JSON 对象: 寻找匹配的 { }
	startIdx := strings.Index(s, "{")
	if startIdx < 0 {
		return s
	}

	depth := 0
	inString := false
	escaped := false
	for i := startIdx; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			continue
		}
		if c == '\\' && inString {
			escaped = true
			continue
		}
		if c == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				return s[startIdx : i+1]
			}
		}
	}

	// fallback: 使用原始的首尾匹配
	endIdx := strings.LastIndex(s, "}")
	if endIdx > startIdx {
		return s[startIdx : endIdx+1]
	}
	return s
}

// fixControlCharsInStrings 修复 JSON 字符串值内的非法控制字符（如字面换行符、制表符等）
func fixControlCharsInStrings(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	inString := false
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			buf.WriteByte(c)
			continue
		}
		if c == '\\' && inString {
			escaped = true
			buf.WriteByte(c)
			continue
		}
		if c == '"' {
			inString = !inString
			buf.WriteByte(c)
			continue
		}
		if inString && c < 0x20 {
			// 替换 JSON 字符串中的非法控制字符
			switch c {
			case '\n':
				buf.WriteString(`\n`)
			case '\r':
				buf.WriteString(`\r`)
			case '\t':
				buf.WriteString(`\t`)
			default:
				buf.WriteString(fmt.Sprintf(`\u%04x`, c))
			}
			continue
		}
		buf.WriteByte(c)
	}
	return buf.String()
}

// fixTrailingCommas 移除 JSON 中 } 或 ] 前的尾部逗号（LLM 常见错误）
func fixTrailingCommas(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	inString := false
	escaped := false
	// pendingComma 缓存遇到的逗号及其后的空白，等确认后续字符再决定是否写入
	pendingComma := ""

	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			buf.WriteByte(c)
			continue
		}
		if c == '\\' && inString {
			escaped = true
			buf.WriteByte(c)
			continue
		}
		if c == '"' {
			inString = !inString
			if pendingComma != "" {
				buf.WriteString(pendingComma)
				pendingComma = ""
			}
			buf.WriteByte(c)
			continue
		}
		if inString {
			buf.WriteByte(c)
			continue
		}
		// 以下为字符串外
		if c == ',' {
			if pendingComma != "" {
				buf.WriteString(pendingComma)
			}
			pendingComma = ","
			continue
		}
		if (c == ' ' || c == '\t' || c == '\n' || c == '\r') && pendingComma != "" {
			pendingComma += string(c)
			continue
		}
		if (c == '}' || c == ']') && pendingComma != "" {
			// 丢弃逗号，保留空白
			if len(pendingComma) > 1 {
				buf.WriteString(pendingComma[1:])
			}
			pendingComma = ""
			buf.WriteByte(c)
			continue
		}
		if pendingComma != "" {
			buf.WriteString(pendingComma)
			pendingComma = ""
		}
		buf.WriteByte(c)
	}
	if pendingComma != "" {
		buf.WriteString(pendingComma)
	}
	return buf.String()
}

func filterOtherLangIDs(word *model.Word, allIDs []string) []string {
	var otherPrefix string
	if word.Language == model.LangEnglish {
		otherPrefix = "ja_"
	} else {
		otherPrefix = "en_"
	}

	var result []string
	for _, id := range allIDs {
		if strings.HasPrefix(id, otherPrefix) {
			result = append(result, id)
		}
	}
	return result
}

func IsAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}
