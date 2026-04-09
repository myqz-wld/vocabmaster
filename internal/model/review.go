package model

import "time"

type Grade int

const (
	GradeAgain Grade = 1
	GradeHard  Grade = 3
	GradeGood  Grade = 4
	GradeEasy  Grade = 5
)

func (g Grade) Label() string {
	switch g {
	case GradeAgain:
		return "Again"
	case GradeHard:
		return "Hard"
	case GradeGood:
		return "Good"
	case GradeEasy:
		return "Easy"
	default:
		return "Unknown"
	}
}

func (g Grade) IsCorrect() bool {
	return g >= GradeHard
}

type ReviewRecord struct {
	WordID       string
	EaseFactor   float64
	Interval     int
	Repetitions  int
	NextReviewAt time.Time
	LastReviewAt time.Time
	TotalReviews int
	CorrectCount int
	FirstSeenAt  time.Time
}

func NewReviewRecord(wordID string, now time.Time) ReviewRecord {
	return ReviewRecord{
		WordID:       wordID,
		EaseFactor:   2.5,
		Interval:     0,
		Repetitions:  0,
		NextReviewAt: now,
		LastReviewAt:  time.Time{},
		TotalReviews: 0,
		CorrectCount: 0,
		FirstSeenAt:  now,
	}
}

type ReviewHistory struct {
	ID               int64
	WordID           string
	Grade            Grade
	ReviewedAt       time.Time
	EaseFactorBefore float64
	EaseFactorAfter  float64
	IntervalBefore   int
	IntervalAfter    int
}
