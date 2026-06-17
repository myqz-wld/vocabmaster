package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

type EnrichResult struct {
	ChineseDef    string          `json:"chinese_def"`
	Pronunciation string          `json:"pronunciation"`
	Examples      []model.Example `json:"examples"`
}

const providerTimeout = 60 * time.Second

type llmProvider struct {
	id         string
	name       string
	executable string
	run        func(context.Context, string) (string, error)
}

var llmProviders = []llmProvider{
	{id: "codex", name: "Codex", executable: "codex", run: runCodex},
	{id: "claude", name: "Claude", executable: "claude", run: runClaude},
}

var execLookPath = exec.LookPath

func EnrichWord(word *model.Word) (*model.Word, error) {
	if word == nil {
		return nil, fmt.Errorf("word is nil")
	}

	prompt := buildEnrichPrompt(word)
	var failures []string
	attempted := false

	for _, provider := range llmProviders {
		if !isProviderAvailable(provider) {
			failures = append(failures, fmt.Sprintf("%s: 未找到 CLI", provider.id))
			continue
		}

		attempted = true
		ctx, cancel := context.WithTimeout(context.Background(), providerTimeout)
		raw, err := provider.run(ctx, prompt)
		if ctx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("调用 %s CLI 超时（60秒）", provider.name)
		}
		cancel()
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", provider.id, err))
			continue
		}

		enriched, err := parseEnrichResponse(word, raw)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", provider.id, err))
			continue
		}
		return enriched, nil
	}

	if !attempted {
		return nil, fmt.Errorf("未找到可用 LLM CLI（尝试顺序: %s）", ProviderOrder())
	}
	return nil, fmt.Errorf("LLM 增强失败（尝试顺序: %s）: %s", ProviderOrder(), strings.Join(failures, "; "))
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

func runCodex(ctx context.Context, prompt string) (string, error) {
	outputFile, err := os.CreateTemp("", "vocabmaster-codex-*.txt")
	if err != nil {
		return "", fmt.Errorf("创建 Codex 输出文件失败: %w", err)
	}
	outputPath := outputFile.Name()
	if err := outputFile.Close(); err != nil {
		_ = os.Remove(outputPath)
		return "", fmt.Errorf("关闭 Codex 输出文件失败: %w", err)
	}
	defer os.Remove(outputPath)

	cmd := exec.CommandContext(
		ctx,
		"codex",
		"exec",
		"--ask-for-approval", "never",
		"--sandbox", "read-only",
		"--skip-git-repo-check",
		"--ephemeral",
		"--ignore-rules",
		"--color", "never",
		"--output-last-message", outputPath,
		prompt,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("调用 Codex CLI 失败: %w%s", err, commandOutputPreview(output))
	}

	result, err := os.ReadFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("读取 Codex 响应失败: %w", err)
	}
	raw := strings.TrimSpace(string(result))
	if raw == "" {
		raw = strings.TrimSpace(string(output))
	}
	if raw == "" {
		return "", fmt.Errorf("Codex 返回了空响应")
	}
	return raw, nil
}

func runClaude(ctx context.Context, prompt string) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", "-p", prompt, "--output-format", "json")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("调用 Claude CLI 失败: %w", err)
	}

	var response struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return "", fmt.Errorf("解析 Claude 响应失败: %w", err)
	}
	if strings.TrimSpace(response.Result) == "" {
		return "", fmt.Errorf("Claude 返回了空响应")
	}
	return response.Result, nil
}

func commandOutputPreview(output []byte) string {
	preview := strings.TrimSpace(string(output))
	if preview == "" {
		return ""
	}
	if len(preview) > 300 {
		preview = preview[:300] + "..."
	}
	return ": " + preview
}

func buildEnrichPrompt(word *model.Word) string {
	var langDesc string
	if word.Language == model.LangEnglish {
		langDesc = "英文"
	} else {
		langDesc = "日文"
	}

	return fmt.Sprintf(`你是一个专业的语言学家。请为以下%s单词提供增强数据。

单词信息：
- 单词: %s
- 语言: %s
- 当前中文释义: %s
- 当前发音: %s
- 词性: %s

请严格按照以下JSON格式返回（不要包含其他文字）：
{
  "chinese_def": "准确、地道的中文释义",
  "pronunciation": "校验后的发音标注（英文用IPA音标，日文用假名）",
  "examples": [
    {"sentence": "自然的例句1", "translation": "中文翻译1"},
    {"sentence": "自然的例句2", "translation": "中文翻译2"}
  ]
}

要求：
1. 中文释义要准确且易懂
2. 不要改变单词、词性或目标语言含义；无法确认时保留当前中文释义或发音
3. 例句必须使用目标语言，中文翻译必须自然准确，例句要常用且不复杂
4. 只返回纯JSON，不要用markdown代码块包裹，不要有多余文字
5. JSON必须严格合法：不要有尾部逗号，字符串内不要有换行`, langDesc, word.Text, word.Language, word.ChineseDef, word.Pronunciation, word.PartOfSpeech)
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

// fixUnescapedQuotes 修复 JSON 字符串值中未转义的双引号
// AI 生成的 JSON 常在中文释义中使用 "引号" 包裹短语，导致 JSON 解析失败
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

		// 找到一个引号，写入并进入字符串内部
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
				// 判断这个引号是否真的是字符串结束符：
				// 跳过空白后，下一个字符应该是 : , } ]
				j := i + 1
				for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n' || s[j] == '\r') {
					j++
				}
				if j >= len(s) || s[j] == ':' || s[j] == ',' || s[j] == '}' || s[j] == ']' {
					buf.WriteByte(c)
					i++
					break
				}
				// 不是结构性引号，转义它
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

func IsAvailable() bool {
	for _, provider := range llmProviders {
		if isProviderAvailable(provider) {
			return true
		}
	}
	return false
}

func ProviderOrder() string {
	names := make([]string, 0, len(llmProviders))
	for _, provider := range llmProviders {
		names = append(names, provider.id)
	}
	return strings.Join(names, " -> ")
}

func isProviderAvailable(provider llmProvider) bool {
	_, err := execLookPath(provider.executable)
	return err == nil
}
