package store

const createTablesSQL = `
CREATE TABLE IF NOT EXISTS review_records (
	word_id       TEXT PRIMARY KEY,
	ease_factor   REAL NOT NULL DEFAULT 2.5,
	interval_days INTEGER NOT NULL DEFAULT 0,
	repetitions   INTEGER NOT NULL DEFAULT 0,
	next_review_at DATETIME,
	last_review_at DATETIME,
	total_reviews  INTEGER NOT NULL DEFAULT 0,
	correct_count  INTEGER NOT NULL DEFAULT 0,
	first_seen_at  DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_review_next ON review_records(next_review_at);

CREATE TABLE IF NOT EXISTS review_history (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	word_id           TEXT NOT NULL,
	grade             INTEGER NOT NULL,
	reviewed_at       DATETIME NOT NULL,
	ease_factor_before REAL,
	ease_factor_after  REAL,
	interval_before   INTEGER,
	interval_after    INTEGER
);

CREATE INDEX IF NOT EXISTS idx_history_word ON review_history(word_id);
CREATE INDEX IF NOT EXISTS idx_history_date ON review_history(reviewed_at);

CREATE TABLE IF NOT EXISTS enriched_words (
	word_id  TEXT PRIMARY KEY,
	data     TEXT NOT NULL
);
`
