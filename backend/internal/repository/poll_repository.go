package repository

import (
	"context"
	"database/sql"
	"flashpoll/internal/domain"
	"time"
)

type PollRepository struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) *PollRepository {
	return &PollRepository{db: db}
}

func (r *PollRepository) GetAll(ctx context.Context) ([]domain.Poll, error) {
	query := `
		SELECT p.id, p.category, p.question, p.created_at,
		       o.id, o.option_text, o.vote_count
		FROM polls p
		LEFT JOIN options o ON p.id = o.poll_id
		ORDER BY p.created_at DESC, o.id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pollMap := make(map[string]*domain.Poll)
	var pollIDs []string

	for rows.Next() {
		var pID, pCategory, pQuestion, pCreatedAt string
		var oID sql.NullInt64
		var oText sql.NullString
		var oVoteCount sql.NullInt32

		if err := rows.Scan(&pID, &pCategory, &pQuestion, &pCreatedAt, &oID, &oText, &oVoteCount); err != nil {
			return nil, err
		}

		poll, exists := pollMap[pID]
		if !exists {
			createdAt, _ := r.parseTime(pCreatedAt)
			poll = &domain.Poll{
				ID:        pID,
				Category:  domain.Category(pCategory),
				Question:  pQuestion,
				CreatedAt: createdAt,
				Options:   []domain.Option{},
			}
			pollMap[pID] = poll
			pollIDs = append(pollIDs, pID)
		}

		if oID.Valid {
			poll.Options = append(poll.Options, domain.Option{
				ID:         oID.Int64,
				PollID:     pID,
				OptionText: oText.String,
				VoteCount:  int(oVoteCount.Int32),
			})
			poll.TotalVotes += int(oVoteCount.Int32)
		}
	}

	result := make([]domain.Poll, 0, len(pollIDs))
	for _, id := range pollIDs {
		poll := pollMap[id]
		// Calculate percentages
		for i := range poll.Options {
			if poll.TotalVotes > 0 {
				poll.Options[i].Percentage = float64(poll.Options[i].VoteCount) / float64(poll.TotalVotes) * 100
			}
		}
		result = append(result, *poll)
	}

	return result, nil
}

func (r *PollRepository) Create(ctx context.Context, poll *domain.Poll) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "INSERT INTO polls (id, category, question, created_at) VALUES (?, ?, ?, ?)",
		poll.ID, poll.Category, poll.Question, poll.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	for i := range poll.Options {
		res, err := tx.ExecContext(ctx, "INSERT INTO options (poll_id, option_text, vote_count) VALUES (?, ?, ?)",
			poll.ID, poll.Options[i].OptionText, 0)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		poll.Options[i].ID = id
		poll.Options[i].PollID = poll.ID
		poll.Options[i].VoteCount = 0
	}

	return tx.Commit()
}

func (r *PollRepository) GetByID(ctx context.Context, id string) (*domain.Poll, error) {
	query := `
		SELECT p.id, p.category, p.question, p.created_at,
		       o.id, o.option_text, o.vote_count
		FROM polls p
		LEFT JOIN options o ON p.id = o.poll_id
		WHERE p.id = ?
		ORDER BY o.id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var poll *domain.Poll
	for rows.Next() {
		var pID, pCategory, pQuestion, pCreatedAt string
		var oID sql.NullInt64
		var oText sql.NullString
		var oVoteCount sql.NullInt32

		if err := rows.Scan(&pID, &pCategory, &pQuestion, &pCreatedAt, &oID, &oText, &oVoteCount); err != nil {
			return nil, err
		}

		if poll == nil {
			createdAt, _ := r.parseTime(pCreatedAt)
			poll = &domain.Poll{
				ID:        pID,
				Category:  domain.Category(pCategory),
				Question:  pQuestion,
				CreatedAt: createdAt,
				Options:   []domain.Option{},
			}
		}

		if oID.Valid {
			poll.Options = append(poll.Options, domain.Option{
				ID:         oID.Int64,
				PollID:     pID,
				OptionText: oText.String,
				VoteCount:  int(oVoteCount.Int32),
			})
			poll.TotalVotes += int(oVoteCount.Int32)
		}
	}

	if poll == nil {
		return nil, sql.ErrNoRows
	}

	// Calculate percentages
	for i := range poll.Options {
		if poll.TotalVotes > 0 {
			poll.Options[i].Percentage = float64(poll.Options[i].VoteCount) / float64(poll.TotalVotes) * 100
		}
	}

	return poll, nil
}

func (r *PollRepository) IncrementVote(ctx context.Context, pollID string, optionID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify poll exists
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(1) FROM polls WHERE id = ?", pollID).Scan(&count)
	if err != nil || count == 0 {
		return sql.ErrNoRows
	}

	// Atomic increment constrained by poll_id
	res, err := tx.ExecContext(ctx, "UPDATE options SET vote_count = vote_count + 1 WHERE id = ? AND poll_id = ?", optionID, pollID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (r *PollRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM polls WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PollRepository) GetStats(ctx context.Context) (int, int, error) {
	var pollCount, voteCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM polls").Scan(&pollCount)
	if err != nil {
		return 0, 0, err
	}
	err = r.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(vote_count), 0) FROM options").Scan(&voteCount)
	if err != nil {
		return 0, 0, err
	}
	return pollCount, voteCount, nil
}
