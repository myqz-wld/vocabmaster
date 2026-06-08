package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

func TestGetReviewCountOnDateUsesLocalCalendarDay(t *testing.T) {
	store, err := NewSQLiteStore(filepath.Join(t.TempDir(), "vocabmaster.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer store.Close()

	loc := time.FixedZone("UTC+8", 8*60*60)
	entries := []model.ReviewHistory{
		{WordID: "en_yesterday", Grade: model.GradeGood, ReviewedAt: time.Date(2026, 6, 7, 23, 30, 0, 0, loc)},
		{WordID: "en_today", Grade: model.GradeGood, ReviewedAt: time.Date(2026, 6, 8, 6, 30, 0, 0, loc)},
	}
	for i := range entries {
		if err := store.AddReviewHistory(&entries[i]); err != nil {
			t.Fatalf("AddReviewHistory(%s) error = %v", entries[i].WordID, err)
		}
	}

	count, err := store.GetReviewCountOnDate(time.Date(2026, 6, 8, 12, 0, 0, 0, loc), "en")
	if err != nil {
		t.Fatalf("GetReviewCountOnDate() error = %v", err)
	}
	if count != 1 {
		t.Fatalf("GetReviewCountOnDate() = %d, want 1", count)
	}
}
