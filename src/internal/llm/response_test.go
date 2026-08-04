package llm

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

func TestCleanJSONResponse(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "纯JSON", input: `{"chinese_def": "测试"}`},
		{name: "代码块在开头", input: "```json\n{\"chinese_def\": \"测试\"}\n```"},
		{name: "代码块在中间有前导文字", input: "Here is the result:\n```json\n{\"chinese_def\": \"测试\"}\n```\nDone."},
		{name: "无语言标记的代码块", input: "```\n{\"chinese_def\": \"测试\"}\n```"},
		{name: "JSON前后有多余文字", input: "以下是结果：\n{\"chinese_def\": \"测试\"}\n希望对你有帮助！"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanJSONResponse(tt.input)
			var value map[string]any
			if err := json.Unmarshal([]byte(got), &value); err != nil {
				t.Errorf("cleanJSONResponse() 结果无法解析为 JSON: %v\n输出: %s", err, got)
			}
		})
	}
}

func TestFixControlCharsInStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "字符串内换行符", input: "{\"def\": \"第一行\n第二行\"}"},
		{name: "字符串内制表符", input: "{\"def\": \"列1\t列2\"}"},
		{name: "字符串内回车换行", input: "{\"def\": \"行1\r\n行2\"}"},
		{name: "字符串外的换行不受影响", input: "{\n  \"def\": \"正常值\"\n}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixControlCharsInStrings(tt.input)
			var value map[string]any
			if err := json.Unmarshal([]byte(got), &value); err != nil {
				t.Errorf("fixControlCharsInStrings() 结果无法解析: %v\n输出: %s", err, got)
			}
		})
	}
}

func TestFixTrailingCommas(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "对象尾部逗号", input: `{"a": "1", "b": "2",}`},
		{name: "数组尾部逗号", input: `{"arr": ["a", "b",]}`},
		{name: "嵌套尾部逗号", input: `{"arr": [{"x": 1,}, {"y": 2,},],}`},
		{name: "逗号后有空白", input: "{\"a\": \"1\",\n  }"},
		{name: "字符串内的逗号+大括号不受影响", input: `{"text": "a,}b"}`},
		{name: "正常JSON不变", input: `{"a": "1", "b": [1, 2]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixTrailingCommas(tt.input)
			var value map[string]any
			if err := json.Unmarshal([]byte(got), &value); err != nil {
				t.Errorf("fixTrailingCommas() 结果无法解析: %v\n输出: %s", err, got)
			}
		})
	}
}

func TestFixUnescapedQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "中文释义中的引号短语", input: `{"chinese_def":"栅栏；界限，范围（常用于 beyond the pale 表示"超出可接受的范围"）"}`},
		{name: "多个未转义引号", input: `{"chinese_def":"表示"你好"和"再见"的用法"}`},
		{name: "正常JSON不变", input: `{"chinese_def":"没有引号的普通文本","pronunciation":"/test/"}`},
		{name: "已转义的引号不受影响", input: `{"chinese_def":"已经转义的\"引号\""}`},
		{name: "值末尾的引号正确处理", input: `{"a":"含"引号"的值","b":"正常值"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixUnescapedQuotes(tt.input)
			var value map[string]any
			if err := json.Unmarshal([]byte(got), &value); err != nil {
				t.Errorf("fixUnescapedQuotes() 结果无法解析: %v\n输入: %s\n输出: %s", err, tt.input, got)
			}
		})
	}
}

func TestFullResponsePipeline(t *testing.T) {
	input := "好的，以下是增强数据：\n```json\n" +
		"{\n" +
		"  \"chinese_def\": \"释义第一行\n第二行\",\n" +
		"  \"pronunciation\": \"/tɛst/\",\n" +
		"  \"examples\": [\n" +
		"    {\"sentence\": \"Example one.\", \"translation\": \"例句一\"},\n" +
		"    {\"sentence\": \"Example two.\", \"translation\": \"例句二\"},\n" +
		"  ],\n" +
		"}\n```\n完成。"

	s := cleanJSONResponse(input)
	s = fixControlCharsInStrings(s)
	s = fixTrailingCommas(s)
	s = fixUnescapedQuotes(s)

	var result EnrichResult
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		t.Fatalf("完整管线解析失败: %v\n输出: %s", err, s)
	}
	if result.ChineseDef == "" || len(result.Examples) != 2 {
		t.Fatalf("unexpected parsed result: %#v", result)
	}
}

func TestBuildEnrichPromptNoLinkedWordIDs(t *testing.T) {
	prompt := buildEnrichPrompt(testWord())
	if strings.Contains(prompt, "linked_word_ids") || strings.Contains(prompt, "关联") {
		t.Error("prompt 不应包含关联词相关内容")
	}
	for _, want := range []string{"test", "不要改变单词、词性或目标语言含义", "例句必须使用目标语言"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt 应包含约束 %q", want)
		}
	}
}

func TestEnrichResultIgnoresLegacyLinkedWordIDs(t *testing.T) {
	input := `{
		"chinese_def": "测试",
		"pronunciation": "/tɛst/",
		"examples": [{"sentence": "A test.", "translation": "一个测试。"}],
		"linked_word_ids": ["ja_tesuto"]
	}`

	var result EnrichResult
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if result.ChineseDef != "测试" || len(result.Examples) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func testWord() *model.Word {
	return &model.Word{
		ID:            "en_test",
		Language:      model.LangEnglish,
		Text:          "test",
		Pronunciation: "/tɛst/",
		ChineseDef:    "测试",
		PartOfSpeech:  "noun",
	}
}
