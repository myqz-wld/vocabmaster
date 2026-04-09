package session

import (
	"fmt"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/sm2"
	"github.com/vocabmaster/vocabmaster/internal/store"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

type ReviewSession struct {
	store  store.Store
	lib    *library.Library
	config Config
}

func NewReviewSession(s store.Store, lib *library.Library, cfg Config) *ReviewSession {
	return &ReviewSession{store: s, lib: lib, config: cfg}
}

func (s *ReviewSession) Run() (*SessionResult, error) {
	now := time.Now()
	// limit: 0 = 全部（不加 LIMIT），>0 = 指定数量，<0 = 回退默认值
	limit := s.config.Count
	if limit < 0 {
		limit = 20
	}

	dueRecords, err := s.store.GetDueWords(now, s.config.Lang, limit)
	if err != nil {
		return nil, fmt.Errorf("获取待复习词汇失败: %w", err)
	}

	if len(dueRecords) == 0 {
		fmt.Println("\n  没有需要复习的单词！")
		return &SessionResult{}, nil
	}

	fmt.Printf("\n  %s开始复习 %d 个单词%s\n", "\033[1m", len(dueRecords), "\033[0m")

	result := &SessionResult{}
	var againWords []model.ReviewRecord

	for i, record := range dueRecords {
		w := s.getWordForDisplay(record.WordID)
		if w == nil {
			continue
		}

		fmt.Printf("\n  [%d/%d]\n", i+1, len(dueRecords))

		var prompt string
		if w.ChineseDef == "" {
			// No Chinese def (e.g. unenriched Japanese words) — show word, recall reading/meaning
			prompt = fmt.Sprintf("单词：%s (%s) → 回忆其含义", w.Text, w.Pronunciation)
		} else if i%2 == 0 {
			prompt = fmt.Sprintf("中文：%s → 回忆对应的单词", w.ChineseDef)
		} else {
			prompt = fmt.Sprintf("单词：%s → 回忆中文释义", w.Text)
		}

		ui.DisplayPrompt(prompt)
		ui.WaitForEnter("按 Enter 查看答案...")

		linkedWords := resolveEnrichedLinkedWords(s.store, s.lib.GetLinkedWordsFor(w))
		ui.DisplayWordCard(w, linkedWords)

		grade := ui.ReadGrade()
		reviewNow := time.Now()

		oldEF := record.EaseFactor
		oldInterval := record.Interval
		updated := sm2.Apply(record, grade, reviewNow)

		if err := s.store.UpsertReviewRecord(&updated); err != nil {
			return nil, err
		}

		history := &model.ReviewHistory{
			WordID:           record.WordID,
			Grade:            grade,
			ReviewedAt:       reviewNow,
			EaseFactorBefore: oldEF,
			EaseFactorAfter:  updated.EaseFactor,
			IntervalBefore:   oldInterval,
			IntervalAfter:    updated.Interval,
		}
		if err := s.store.AddReviewHistory(history); err != nil {
			return nil, err
		}

		result.Total++
		result.Reviewed++
		if grade.IsCorrect() {
			result.Correct++
		} else {
			againWords = append(againWords, updated)
		}
	}

	// Re-review words marked "Again"
	if len(againWords) > 0 {
		fmt.Printf("\n  %s--- 再次复习 %d 个错误词 ---%s\n", "\033[1m", len(againWords), "\033[0m")
		for i, record := range againWords {
			w := s.getWordForDisplay(record.WordID)
			if w == nil {
				continue
			}

			fmt.Printf("\n  [再复习 %d/%d]\n", i+1, len(againWords))
			ui.DisplayPrompt(fmt.Sprintf("中文：%s → 回忆对应的单词", w.ChineseDef))
			ui.WaitForEnter("按 Enter 查看答案...")
			ui.DisplayReveal(w)

			grade := ui.ReadGrade()
			reviewNow := time.Now()

			current, err := s.store.GetReviewRecord(record.WordID)
			if err != nil {
				return nil, err
			}
			if current == nil {
				continue
			}

			oldEF := current.EaseFactor
			oldInterval := current.Interval
			updated := sm2.Apply(*current, grade, reviewNow)
			if err := s.store.UpsertReviewRecord(&updated); err != nil {
				return nil, err
			}

			history := &model.ReviewHistory{
				WordID:           record.WordID,
				Grade:            grade,
				ReviewedAt:       reviewNow,
				EaseFactorBefore: oldEF,
				EaseFactorAfter:  updated.EaseFactor,
				IntervalBefore:   oldInterval,
				IntervalAfter:    updated.Interval,
			}
			if err := s.store.AddReviewHistory(history); err != nil {
				return nil, err
			}

			result.Total++
			if grade.IsCorrect() {
				result.Correct++
			}
		}
	}

	return result, nil
}

func (s *ReviewSession) getWordForDisplay(wordID string) *model.Word {
	cached, err := s.store.GetEnrichedWord(wordID)
	if err == nil && cached != nil {
		return cached
	}
	return s.lib.GetWord(wordID)
}
