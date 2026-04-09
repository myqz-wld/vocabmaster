package session

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/library"
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
		rand.Shuffle(len(newWords), func(i, j int) {
			newWords[i], newWords[j] = newWords[j], newWords[i]
		})
		newWords = newWords[:count]
	}

	fmt.Printf("\n  %s开始学习 %d 个新词%s\n", "\033[1m", len(newWords), "\033[0m")

	// Phase 1: Introduction
	for i, w := range newWords {
		fmt.Printf("\n  [%d/%d]\n", i+1, len(newWords))

		displayWord := enrichWord(s.store, w)
		ui.DisplayWordCard(displayWord)

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
		displayWord := enrichWord(s.store, w)

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
