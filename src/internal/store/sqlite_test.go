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

func TestGetReviewCountOnDateNormalizesReviewedAtOffsets(t *testing.T) {
	store, err := NewSQLiteStore(filepath.Join(t.TempDir(), "vocabmaster.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer store.Close()

	local := time.FixedZone("UTC+8", 8*60*60)
	west := time.FixedZone("UTC-7", -7*60*60)
	farEast := time.FixedZone("UTC+14", 14*60*60)
	entries := []model.ReviewHistory{
		// 2026-06-08 14:30 in UTC+8, but the stored text starts with 2026-06-07.
		{WordID: "en_inside", Grade: model.GradeGood, ReviewedAt: time.Date(2026, 6, 7, 23, 30, 0, 0, west)},
		// 2026-06-07 18:30 in UTC+8, but the stored text starts with 2026-06-08.
		{WordID: "ja_outside", Grade: model.GradeGood, ReviewedAt: time.Date(2026, 6, 8, 0, 30, 0, 0, farEast)},
	}
	for i := range entries {
		if err := store.AddReviewHistory(&entries[i]); err != nil {
			t.Fatalf("AddReviewHistory(%s) error = %v", entries[i].WordID, err)
		}
	}

	date := time.Date(2026, 6, 8, 12, 0, 0, 0, local)
	enCount, err := store.GetReviewCountOnDate(date, "en")
	if err != nil {
		t.Fatalf("GetReviewCountOnDate(en) error = %v", err)
	}
	if enCount != 1 {
		t.Fatalf("GetReviewCountOnDate(en) = %d, want 1", enCount)
	}

	jaCount, err := store.GetReviewCountOnDate(date, "ja")
	if err != nil {
		t.Fatalf("GetReviewCountOnDate(ja) error = %v", err)
	}
	if jaCount != 0 {
		t.Fatalf("GetReviewCountOnDate(ja) = %d, want 0", jaCount)
	}
}
