package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/model"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}

	if _, err := db.Exec(createTablesSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) GetReviewRecord(wordID string) (*model.ReviewRecord, error) {
	row := s.db.QueryRow(`SELECT word_id, ease_factor, interval_days, repetitions,
		next_review_at, last_review_at, total_reviews, correct_count, first_seen_at
		FROM review_records WHERE word_id = ?`, wordID)

	var r model.ReviewRecord
	var nextReview, lastReview, firstSeen string
	err := row.Scan(&r.WordID, &r.EaseFactor, &r.Interval, &r.Repetitions,
		&nextReview, &lastReview, &r.TotalReviews, &r.CorrectCount, &firstSeen)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	r.NextReviewAt = parseTime(nextReview)
	r.LastReviewAt = parseTime(lastReview)
	r.FirstSeenAt = parseTime(firstSeen)
	return &r, nil
}

func (s *SQLiteStore) UpsertReviewRecord(r *model.ReviewRecord) error {
	_, err := s.db.Exec(`INSERT INTO review_records
		(word_id, ease_factor, interval_days, repetitions, next_review_at, last_review_at, total_reviews, correct_count, first_seen_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(word_id) DO UPDATE SET
			ease_factor = excluded.ease_factor,
			interval_days = excluded.interval_days,
			repetitions = excluded.repetitions,
			next_review_at = excluded.next_review_at,
			last_review_at = excluded.last_review_at,
			total_reviews = excluded.total_reviews,
			correct_count = excluded.correct_count`,
		r.WordID, r.EaseFactor, r.Interval, r.Repetitions,
		formatTime(r.NextReviewAt), formatTime(r.LastReviewAt),
		r.TotalReviews, r.CorrectCount, formatTime(r.FirstSeenAt))
	return err
}

func (s *SQLiteStore) ResetWord(wordID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM review_records WHERE word_id = ?`, wordID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM review_history WHERE word_id = ?`, wordID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLiteStore) ResetAll() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM review_records`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM review_history`); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLiteStore) GetDueWords(now time.Time, lang string, limit int) ([]model.ReviewRecord, error) {
	query := `SELECT word_id, ease_factor, interval_days, repetitions,
		next_review_at, last_review_at, total_reviews, correct_count, first_seen_at
		FROM review_records WHERE next_review_at <= ?`
	args := []any{formatTime(now)}

	if lang != "" && lang != "all" {
		query += ` AND word_id LIKE ?`
		args = append(args, lang+"_%")
	}

	query += ` ORDER BY next_review_at ASC`
	if limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}

	return s.queryRecords(query, args...)
}

func (s *SQLiteStore) GetLearnedWordIDs() ([]string, error) {
	rows, err := s.db.Query(`SELECT word_id FROM review_records`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *SQLiteStore) GetLearnedWordSet() (map[string]bool, error) {
	ids, err := s.GetLearnedWordIDs()
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set, nil
}

func (s *SQLiteStore) GetAllReviewRecords() ([]model.ReviewRecord, error) {
	return s.queryRecords(`SELECT word_id, ease_factor, interval_days, repetitions,
		next_review_at, last_review_at, total_reviews, correct_count, first_seen_at
		FROM review_records`)
}

func (s *SQLiteStore) AddReviewHistory(entry *model.ReviewHistory) error {
	_, err := s.db.Exec(`INSERT INTO review_history
		(word_id, grade, reviewed_at, ease_factor_before, ease_factor_after, interval_before, interval_after)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entry.WordID, int(entry.Grade), formatTime(entry.ReviewedAt),
		entry.EaseFactorBefore, entry.EaseFactorAfter, entry.IntervalBefore, entry.IntervalAfter)
	return err
}

func (s *SQLiteStore) GetEnrichedWord(wordID string) (*model.Word, error) {
	var data string
	err := s.db.QueryRow(`SELECT data FROM enriched_words WHERE word_id = ?`, wordID).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var w model.Word
	if err := json.Unmarshal([]byte(data), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *SQLiteStore) SaveEnrichedWord(word *model.Word) error {
	data, err := json.Marshal(word)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT OR REPLACE INTO enriched_words (word_id, data) VALUES (?, ?)`,
		word.ID, string(data))
	return err
}

func (s *SQLiteStore) GetStats(lang string) (*Stats, error) {
	stats := &Stats{}

	whereClause := ""
	var args []any
	if lang != "" && lang != "all" {
		whereClause = ` WHERE word_id LIKE ?`
		args = []any{lang + "_%"}
	}

	err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total_reviews), 0), COALESCE(SUM(correct_count), 0)
		FROM review_records`+whereClause, args...).Scan(&stats.LearnedWords, &stats.TotalReviews, &stats.CorrectCount)
	if err != nil {
		return nil, err
	}

	masteredArgs := args
	masteredWhere := whereClause
	if masteredWhere == "" {
		masteredWhere = ` WHERE ease_factor >= 2.5 AND interval_days >= 21`
	} else {
		masteredWhere += ` AND ease_factor >= 2.5 AND interval_days >= 21`
	}
	err = s.db.QueryRow(`SELECT COUNT(*) FROM review_records`+masteredWhere, masteredArgs...).Scan(&stats.MasteredWords)
	if err != nil {
		return nil, err
	}

	now := formatTime(time.Now())
	dueArgs := []any{now}
	dueWhere := ` WHERE next_review_at <= ?`
	if lang != "" && lang != "all" {
		dueWhere += ` AND word_id LIKE ?`
		dueArgs = append(dueArgs, lang+"_%")
	}
	err = s.db.QueryRow(`SELECT COUNT(*) FROM review_records`+dueWhere, dueArgs...).Scan(&stats.DueWords)
	if err != nil {
		return nil, err
	}

	stats.Streak = s.calculateStreak(lang)
	return stats, nil
}

func (s *SQLiteStore) calculateStreak(lang string) int {
	streak := 0
	today := time.Now().Truncate(24 * time.Hour)

	// 如果今天还没有学习记录，从昨天开始回溯
	start := today
	count, err := s.GetReviewCountOnDate(start, lang)
	if err != nil || count == 0 {
		start = start.AddDate(0, 0, -1)
	}

	for i := 0; ; i++ {
		date := start.AddDate(0, 0, -i)
		count, err := s.GetReviewCountOnDate(date, lang)
		if err != nil || count == 0 {
			break
		}
		streak++
	}
	return streak
}

func (s *SQLiteStore) GetReviewCountOnDate(date time.Time, lang string) (int, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 1)

	query := `SELECT COUNT(*) FROM review_history WHERE reviewed_at >= ? AND reviewed_at < ?`
	args := []any{formatTime(start), formatTime(end)}

	if lang != "" && lang != "all" {
		query += ` AND word_id LIKE ?`
		args = append(args, lang+"_%")
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (s *SQLiteStore) queryRecords(query string, args ...any) ([]model.ReviewRecord, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.ReviewRecord
	for rows.Next() {
		var r model.ReviewRecord
		var nextReview, lastReview, firstSeen string
		if err := rows.Scan(&r.WordID, &r.EaseFactor, &r.Interval, &r.Repetitions,
			&nextReview, &lastReview, &r.TotalReviews, &r.CorrectCount, &firstSeen); err != nil {
			return nil, err
		}
		r.NextReviewAt = parseTime(nextReview)
		r.LastReviewAt = parseTime(lastReview)
		r.FirstSeenAt = parseTime(firstSeen)
		records = append(records, r)
	}
	return records, rows.Err()
}

func parseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
