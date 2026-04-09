package store

import (
	"time"

	"github.com/vocabmaster/vocabmaster/internal/model"
)

type Stats struct {
	TotalWords    int
	LearnedWords  int
	MasteredWords int
	DueWords      int
	TotalReviews  int
	CorrectCount  int
	Streak        int
}

type Store interface {
	GetReviewRecord(wordID string) (*model.ReviewRecord, error)
	UpsertReviewRecord(record *model.ReviewRecord) error
	GetDueWords(now time.Time, lang string, limit int) ([]model.ReviewRecord, error)
	GetLearnedWordIDs() ([]string, error)
	GetLearnedWordSet() (map[string]bool, error)
	GetAllReviewRecords() ([]model.ReviewRecord, error)

	AddReviewHistory(entry *model.ReviewHistory) error

	ResetWord(wordID string) error
	ResetAll() error

	GetEnrichedWord(wordID string) (*model.Word, error)
	SaveEnrichedWord(word *model.Word) error

	GetStats(lang string) (*Stats, error)
	GetReviewCountOnDate(date time.Time, lang string) (int, error)

	Close() error
}
