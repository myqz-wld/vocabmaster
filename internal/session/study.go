package session

import (
	"fmt"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/store"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

type StudySession struct {
	store  store.Store
	lib    *library.Library
	config Config
}

func NewStudySession(s store.Store, lib *library.Library, cfg Config) *StudySession {
	return &StudySession{store: s, lib: lib, config: cfg}
}

func (s *StudySession) Run() error {
	now := time.Now()

	dueRecords, err := s.store.GetDueWords(now, s.config.Lang, 0)
	if err != nil {
		return fmt.Errorf("查询到期词汇失败: %w", err)
	}

	dueCount := len(dueRecords)
	totalResult := &SessionResult{}

	// Phase 1: Review due words
	if dueCount > 0 {
		reviewCount := dueCount
		if reviewCount > 30 {
			reviewCount = 30
		}

		reviewCfg := Config{
			Lang:  s.config.Lang,
			Count: reviewCount,
		}

		reviewSession := NewReviewSession(s.store, s.lib, reviewCfg)
		reviewResult, err := reviewSession.Run()
		if err != nil {
			return err
		}

		totalResult.Reviewed += reviewResult.Reviewed
		totalResult.Correct += reviewResult.Correct
		totalResult.Total += reviewResult.Total
	}

	// Phase 2: Learn new words based on load
	newWordCount := calculateNewWords(dueCount, s.config.NewWords)

	if newWordCount > 0 {
		learnedSet, err := s.store.GetLearnedWordSet()
		if err != nil {
			return err
		}

		available := s.lib.GetUnlearnedWords(s.config.Lang, model.DifficultyLevel(s.config.Level), learnedSet)
		if len(available) > 0 {
			if newWordCount > len(available) {
				newWordCount = len(available)
			}

			fmt.Printf("\n  %s--- 进入新词学习阶段 (%d 词) ---%s\n", "\033[1m", newWordCount, "\033[0m")

			learnCfg := Config{
				Lang:  s.config.Lang,
				Level: s.config.Level,
				Count: newWordCount,
			}

			learnSession := NewLearnSession(s.store, s.lib, learnCfg)
			learnResult, err := learnSession.Run()
			if err != nil {
				return err
			}

			totalResult.Learned += learnResult.Learned
			totalResult.Correct += learnResult.Correct
			totalResult.Total += learnResult.Total
		}
	}

	ui.DisplaySessionSummary(totalResult.Reviewed, totalResult.Learned, totalResult.Correct, totalResult.Total)
	return nil
}

func calculateNewWords(dueCount, override int) int {
	if override > 0 {
		return override
	}

	const dailyTarget = 10

	switch {
	case dueCount > 20:
		return 0
	case dueCount > 10:
		return dailyTarget / 2
	default:
		return dailyTarget
	}
}
