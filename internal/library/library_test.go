package library

import (
	"testing"

	"github.com/vocabmaster/vocabmaster/internal/model"
)

func TestNewLibrary(t *testing.T) {
	lib, err := NewLibrary()
	if err != nil {
		t.Fatalf("NewLibrary() 失败: %v", err)
	}
	if lib.TotalCount() == 0 {
		t.Error("词库不应为空")
	}
}

func TestFilterByLang(t *testing.T) {
	lib, _ := NewLibrary()

	en := lib.FilterByLang("en")
	ja := lib.FilterByLang("ja")
	all := lib.FilterByLang("all")

	if len(en) == 0 {
		t.Error("英语词库不应为空")
	}
	if len(ja) == 0 {
		t.Error("日语词库不应为空")
	}
	if len(all) != len(en)+len(ja) {
		t.Errorf("all(%d) != en(%d) + ja(%d)", len(all), len(en), len(ja))
	}
}

func TestGetUnlearnedWords(t *testing.T) {
	lib, _ := NewLibrary()

	learned := map[string]bool{}
	words := lib.GetUnlearnedWords("en", model.Beginner, learned)
	if len(words) == 0 {
		t.Error("未学词汇不应为空")
	}

	// 标记所有返回词汇为已学习，应该不再出现
	for _, w := range words {
		learned[w.ID] = true
	}
	remaining := lib.GetUnlearnedWords("en", model.Beginner, learned)
	if len(remaining) != 0 {
		t.Errorf("标记全部已学后仍有 %d 个未学词汇", len(remaining))
	}
}

func TestSearch(t *testing.T) {
	lib, _ := NewLibrary()

	results := lib.Search("hello")
	if len(results) == 0 {
		t.Error("搜索 'hello' 应有结果")
	}

	results = lib.Search("xyznonexistent12345")
	if len(results) != 0 {
		t.Error("搜索不存在的词应返回空")
	}
}
