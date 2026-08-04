package llm

import (
	"fmt"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

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
