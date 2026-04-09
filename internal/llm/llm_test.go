package llm

import (
	"encoding/json"
	"testing"
)

func TestCleanJSONResponse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string // 期望能被 json.Unmarshal 解析
	}{
		{
			name:  "纯JSON",
			input: `{"chinese_def": "测试"}`,
			want:  `{"chinese_def": "测试"}`,
		},
		{
			name:  "代码块在开头",
			input: "```json\n{\"chinese_def\": \"测试\"}\n```",
			want:  `{"chinese_def": "测试"}`,
		},
		{
			name:  "代码块在中间有前导文字",
			input: "Here is the result:\n```json\n{\"chinese_def\": \"测试\"}\n```\nDone.",
			want:  `{"chinese_def": "测试"}`,
		},
		{
			name:  "无语言标记的代码块",
			input: "```\n{\"chinese_def\": \"测试\"}\n```",
			want:  `{"chinese_def": "测试"}`,
		},
		{
			name:  "JSON前后有多余文字",
			input: "以下是结果：\n{\"chinese_def\": \"测试\"}\n希望对你有帮助！",
			want:  `{"chinese_def": "测试"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanJSONResponse(tt.input)
			var v map[string]interface{}
			if err := json.Unmarshal([]byte(got), &v); err != nil {
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
		{
			name:  "字符串内换行符",
			input: "{\"def\": \"第一行\n第二行\"}",
		},
		{
			name:  "字符串内制表符",
			input: "{\"def\": \"列1\t列2\"}",
		},
		{
			name:  "字符串内回车换行",
			input: "{\"def\": \"行1\r\n行2\"}",
		},
		{
			name:  "字符串外的换行不受影响",
			input: "{\n  \"def\": \"正常值\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixControlCharsInStrings(tt.input)
			var v map[string]interface{}
			if err := json.Unmarshal([]byte(got), &v); err != nil {
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
		{
			name:  "对象尾部逗号",
			input: `{"a": "1", "b": "2",}`,
		},
		{
			name:  "数组尾部逗号",
			input: `{"arr": ["a", "b",]}`,
		},
		{
			name:  "嵌套尾部逗号",
			input: `{"arr": [{"x": 1,}, {"y": 2,},],}`,
		},
		{
			name:  "逗号后有空白",
			input: "{\"a\": \"1\",\n  }",
		},
		{
			name:  "字符串内的逗号+大括号不受影响",
			input: `{"text": "a,}b"}`,
		},
		{
			name:  "正常JSON不变",
			input: `{"a": "1", "b": [1, 2]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixTrailingCommas(tt.input)
			var v map[string]interface{}
			if err := json.Unmarshal([]byte(got), &v); err != nil {
				t.Errorf("fixTrailingCommas() 结果无法解析: %v\n输出: %s", err, got)
			}
		})
	}
}

func TestFullPipeline(t *testing.T) {
	// 模拟一个典型的 LLM 返回，包含多种问题
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

	var result EnrichResult
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		t.Fatalf("完整管线解析失败: %v\n输出: %s", err, s)
	}
	if result.ChineseDef == "" {
		t.Error("chinese_def 不应为空")
	}
	if len(result.Examples) != 2 {
		t.Errorf("期望 2 个例句，得到 %d", len(result.Examples))
	}
}
