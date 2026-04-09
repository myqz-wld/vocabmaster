package session

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/llm"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/sm2"
	"github.com/vocabmaster/vocabmaster/internal/store"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

type LearnSession struct {
	store  store.Store
	lib    *library.Library
	config Config
}

func NewLearnSession(s store.Store, lib *library.Library, cfg Config) *LearnSession {
	return &LearnSession{store: s, lib: lib, config: cfg}
}

func (s *LearnSession) Run() (*SessionResult, error) {
	learnedSet, err := s.store.GetLearnedWordSet()
	if err != nil {
		return nil, fmt.Errorf("获取已学词汇失败: %w", err)
	}

	newWords := s.lib.GetUnlearnedWords(s.config.Lang, s.config.Level, learnedSet)
	count := s.config.Count
	if count <= 0 {
		count = 5
	}
	if len(newWords) == 0 {
		fmt.Println("\n  没有新词可学！所有词汇都已学习过。")
		return &SessionResult{}, nil
	}
	if len(newWords) > count {
		newWords = newWords[:count]
	}

	fmt.Printf("\n  %s开始学习 %d 个新词%s\n", "\033[1m", len(newWords), "\033[0m")

	// Phase 1: Introduction
	allWordIDs := getAllWordIDs(s.lib)
	for i, w := range newWords {
		fmt.Printf("\n  [%d/%d]\n", i+1, len(newWords))

		displayWord := s.enrichWord(w, allWordIDs)
		linkedWords := resolveEnrichedLinkedWords(s.store, s.lib.GetLinkedWordsFor(displayWord))
		ui.DisplayWordCard(displayWord, linkedWords)

		now := time.Now()
		record := model.NewReviewRecord(w.ID, now)
		if err := s.store.UpsertReviewRecord(&record); err != nil {
			return nil, fmt.Errorf("保存学习记录失败: %w", err)
		}

		ui.WaitForEnter("按 Enter 继续下一个词...")
	}

	// Phase 2: Quick quiz
	fmt.Printf("\n  %s--- 快速测验 ---%s\n", "\033[1m", "\033[0m")

	shuffled := make([]*model.Word, len(newWords))
	copy(shuffled, newWords)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := &SessionResult{Learned: len(newWords)}
	for _, w := range shuffled {
		displayWord := s.getDisplayWord(w)

		var prompt string
		if w.Language == model.LangEnglish {
			prompt = fmt.Sprintf("中文：%s → 请回忆对应的英文单词", displayWord.ChineseDef)
		} else {
			prompt = fmt.Sprintf("中文：%s → 请回忆对应的日文单词", displayWord.ChineseDef)
		}
		ui.DisplayPrompt(prompt)
		ui.WaitForEnter("按 Enter 查看答案...")
		ui.DisplayReveal(displayWord)

		grade := ui.ReadGrade()
		now := time.Now()

		record, err := s.store.GetReviewRecord(w.ID)
		if err != nil {
			return nil, err
		}
		if record == nil {
			r := model.NewReviewRecord(w.ID, now)
			record = &r
		}

		oldEF := record.EaseFactor
		oldInterval := record.Interval
		*record = sm2.Apply(*record, grade, now)

		if err := s.store.UpsertReviewRecord(record); err != nil {
			return nil, err
		}

		history := &model.ReviewHistory{
			WordID:           w.ID,
			Grade:            grade,
			ReviewedAt:       now,
			EaseFactorBefore: oldEF,
			EaseFactorAfter:  record.EaseFactor,
			IntervalBefore:   oldInterval,
			IntervalAfter:    record.Interval,
		}
		if err := s.store.AddReviewHistory(history); err != nil {
			return nil, err
		}

		result.Total++
		if grade.IsCorrect() {
			result.Correct++
		}
	}

	return result, nil
}

func (s *LearnSession) enrichWord(w *model.Word, allWordIDs []string) *model.Word {
	cached, err := s.store.GetEnrichedWord(w.ID)
	if err == nil && cached != nil {
		return cached
	}

	if llm.IsAvailable() {
		fmt.Printf("  %s正在通过 AI 增强词汇数据...%s\n", "\033[90m", "\033[0m")
		enriched, err := llm.EnrichWord(w, allWordIDs)
		if err == nil {
			if saveErr := s.store.SaveEnrichedWord(enriched); saveErr != nil {
				fmt.Printf("  %s保存增强数据失败: %v%s\n", "\033[90m", saveErr, "\033[0m")
			}
			return enriched
		}
		fmt.Printf("  %sAI 增强失败，使用基础数据: %v%s\n", "\033[90m", err, "\033[0m")
	}

	return w
}

func (s *LearnSession) getDisplayWord(w *model.Word) *model.Word {
	cached, err := s.store.GetEnrichedWord(w.ID)
	if err == nil && cached != nil {
		return cached
	}
	return w
}
