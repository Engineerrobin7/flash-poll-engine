package db

import "database/sql"

const Schema = `
CREATE TABLE IF NOT EXISTS polls (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL CHECK (category IN ('Tech', 'Business', 'Design')),
    question TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    poll_id TEXT NOT NULL,
    option_text TEXT NOT NULL,
    vote_count INTEGER NOT NULL DEFAULT 0 CHECK (vote_count >= 0),
    FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
    UNIQUE (poll_id, option_text)
);

CREATE INDEX IF NOT EXISTS idx_polls_created_at_desc ON polls(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_options_poll_id ON options(poll_id);
`

func Migrate(db *sql.DB) error {
	_, err := db.Exec(Schema)
	return err
}
