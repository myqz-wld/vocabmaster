package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

type EnrichResult struct {
	ChineseDef    string          `json:"chinese_def"`
	Pronunciation string          `json:"pronunciation"`
	Examples      []model.Example `json:"examples"`
}

func parseEnrichResponse(word *model.Word, raw string) (*model.Word, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("AI 返回了空响应")
	}

	resultJSON := cleanJSONResponse(raw)
	resultJSON = fixControlCharsInStrings(resultJSON)
	resultJSON = fixTrailingCommas(resultJSON)
	resultJSON = fixUnescapedQuotes(resultJSON)

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
		for i := range result.Examples {
			result.Examples[i].Sentence = strings.Join(strings.Fields(result.Examples[i].Sentence), " ")
			result.Examples[i].Translation = strings.Join(strings.Fields(result.Examples[i].Translation), " ")
		}
		enriched.Examples = result.Examples
	}
	return &enriched, nil
}

func cleanJSONResponse(raw string) string {
	s := strings.TrimSpace(raw)

	// 移除 markdown 代码块标记（支持代码块出现在任意位置）
	if idx := strings.Index(s, "```"); idx >= 0 {
		afterOpen := s[idx+3:]
		nlIdx := strings.Index(afterOpen, "\n")
		if nlIdx >= 0 {
			inner := afterOpen[nlIdx+1:]
			if closeIdx := strings.Index(inner, "```"); closeIdx >= 0 {
				s = strings.TrimSpace(inner[:closeIdx])
			} else {
				s = strings.TrimSpace(inner)
			}
		}
	}

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

	endIdx := strings.LastIndex(s, "}")
	if endIdx > startIdx {
		return s[startIdx : endIdx+1]
	}
	return s
}

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

func fixTrailingCommas(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	inString := false
	escaped := false
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

func fixUnescapedQuotes(s string) string {
	var buf strings.Builder
	buf.Grow(len(s) + 20)

	i := 0
	for i < len(s) {
		c := s[i]
		if c != '"' {
			buf.WriteByte(c)
			i++
			continue
		}

		buf.WriteByte(c)
		i++
		for i < len(s) {
			c = s[i]
			if c == '\\' {
				buf.WriteByte(c)
				i++
				if i < len(s) {
					buf.WriteByte(s[i])
					i++
				}
				continue
			}
			if c == '"' {
				j := i + 1
				for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n' || s[j] == '\r') {
					j++
				}
				if j >= len(s) || s[j] == ':' || s[j] == ',' || s[j] == '}' || s[j] == ']' {
					buf.WriteByte(c)
					i++
					break
				}
				buf.WriteString(`\"`)
				i++
				continue
			}
			buf.WriteByte(c)
			i++
		}
	}
	return buf.String()
}
