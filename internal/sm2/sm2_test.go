package sm2

import (
	"testing"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/model"
)

func TestFirstReviewEasy(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	r = Apply(r, model.GradeEasy, now)

	if r.Interval != 1 {
		t.Errorf("expected interval 1, got %d", r.Interval)
	}
	if r.Repetitions != 1 {
		t.Errorf("expected repetitions 1, got %d", r.Repetitions)
	}
}

func TestSecondReviewGood(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	r = Apply(r, model.GradeGood, now)
	r = Apply(r, model.GradeGood, now.AddDate(0, 0, 1))

	if r.Interval != 6 {
		t.Errorf("expected interval 6, got %d", r.Interval)
	}
	if r.Repetitions != 2 {
		t.Errorf("expected repetitions 2, got %d", r.Repetitions)
	}
}

func TestThirdReviewUsesEF(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	// 3x GradeGood: EF stays 2.5, third interval = round(6 * 2.5) = 15
	r = Apply(r, model.GradeGood, now)
	r = Apply(r, model.GradeGood, now.AddDate(0, 0, 1))
	r = Apply(r, model.GradeGood, now.AddDate(0, 0, 7))

	if r.Interval != 15 {
		t.Errorf("expected interval 15, got %d", r.Interval)
	}
	if r.Repetitions != 3 {
		t.Errorf("expected repetitions 3, got %d", r.Repetitions)
	}
}

func TestThirdReviewUsesNewEFNotOld(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	// 2x GradeGood → EF=2.5, interval=6
	// 3rd GradeHard → newEF=2.36, interval should use newEF: round(6 * 2.36) = 14
	// If it used oldEF (2.5): round(6 * 2.5) = 15 — wrong
	r = Apply(r, model.GradeGood, now)
	r = Apply(r, model.GradeGood, now.AddDate(0, 0, 1))
	r = Apply(r, model.GradeHard, now.AddDate(0, 0, 7))

	if r.Interval != 14 {
		t.Errorf("expected interval 14 (using newEF=2.36), got %d (would be 15 if using oldEF)", r.Interval)
	}
}

func TestIncorrectResetsRepetitions(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	r = Apply(r, model.GradeGood, now)
	r = Apply(r, model.GradeGood, now.AddDate(0, 0, 1))
	r = Apply(r, model.GradeAgain, now.AddDate(0, 0, 7))

	if r.Repetitions != 0 {
		t.Errorf("expected repetitions 0, got %d", r.Repetitions)
	}
	if r.Interval != 1 {
		t.Errorf("expected interval 1, got %d", r.Interval)
	}
}

func TestEFNeverBelowMinimum(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	for i := 0; i < 20; i++ {
		r = Apply(r, model.GradeAgain, now.AddDate(0, 0, i))
	}

	if r.EaseFactor < 1.3 {
		t.Errorf("EF should never be below 1.3, got %f", r.EaseFactor)
	}
}

func TestEFDecreasesOnIncorrect(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)
	initialEF := r.EaseFactor

	r = Apply(r, model.GradeAgain, now)

	if r.EaseFactor >= initialEF {
		t.Errorf("EF should decrease after incorrect answer, got %f (initial %f)", r.EaseFactor, initialEF)
	}
}

func TestEFConvergesToMinimum(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	for i := 0; i < 50; i++ {
		r = Apply(r, model.GradeAgain, now.AddDate(0, 0, i))
	}

	if r.EaseFactor != 1.3 {
		t.Errorf("EF should converge to 1.3 after many incorrect answers, got %f", r.EaseFactor)
	}
}

func TestTotalReviewsAndCorrectCount(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.NewReviewRecord("test_word", now)

	r = Apply(r, model.GradeGood, now)
	r = Apply(r, model.GradeAgain, now.AddDate(0, 0, 1))
	r = Apply(r, model.GradeEasy, now.AddDate(0, 0, 2))

	if r.TotalReviews != 3 {
		t.Errorf("expected 3 total reviews, got %d", r.TotalReviews)
	}
	if r.CorrectCount != 2 {
		t.Errorf("expected 2 correct, got %d", r.CorrectCount)
	}
}
