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

	resultJSON := response.Result
	startIdx := strings.Index(resultJSON, "{")
	endIdx := strings.LastIndex(resultJSON, "}")
	if startIdx >= 0 && endIdx > startIdx {
		resultJSON = resultJSON[startIdx : endIdx+1]
	}

	var result EnrichResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, fmt.Errorf("解析增强数据失败: %w", err)
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
4. 如果找不到合适的关联词，linked_word_ids 返回空数组`, langDesc, word.Text, word.Language, word.ChineseDef, word.Pronunciation, word.PartOfSpeech, strings.Join(otherLangIDs, ", "))
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
