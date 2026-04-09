package session

import (
	"fmt"

	"github.com/vocabmaster/vocabmaster/internal/llm"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/store"
)

type Config struct {
	Lang     string
	Level    model.DifficultyLevel
	Count    int
	NewWords int
}

type SessionResult struct {
	Reviewed int
	Learned  int
	Correct  int
	Total    int
}

// enrichWord 获取词汇的 AI 增强数据，优先使用缓存，无缓存时调用 AI 生成
func enrichWord(s store.Store, w *model.Word) *model.Word {
	cached, err := s.GetEnrichedWord(w.ID)
	if err == nil && cached != nil {
		return cached
	}

	if llm.IsAvailable() {
		fmt.Printf("  %s正在通过 AI 增强词汇数据...%s\n", "\033[90m", "\033[0m")
		enriched, err := llm.EnrichWord(w)
		if err == nil {
			if saveErr := s.SaveEnrichedWord(enriched); saveErr != nil {
				fmt.Printf("  %s保存增强数据失败: %v%s\n", "\033[90m", saveErr, "\033[0m")
			}
			return enriched
		}
		fmt.Printf("  %sAI 增强失败，使用基础数据: %v%s\n", "\033[90m", err, "\033[0m")
	}

	return w
}
