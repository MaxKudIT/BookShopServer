package reading

import (
	"context"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (rs *readingStorage) Start(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, sessionId uuid.UUID, startedAt time.Time) (domain.ReadingState, error) {
	tx, err := rs.db.BeginTx(ctx, nil)
	if err != nil {
		rs.l.Error("Transaction begin failed", "error", err)
		return domain.ReadingState{}, err
	}
	defer tx.Rollback()

	readingState := domain.ReadingState{BookId: bookId}
	const UpsertReadingProgressQuery = `
		INSERT INTO user_reading_progress (user_id, book_id, status, current_page, progress_percent, updated_at)
		VALUES ($1, $2, $3, 1, 0, NOW())
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET
			status = CASE
				WHEN user_reading_progress.status = $4 THEN user_reading_progress.status
				ELSE $3
			END,
			updated_at = NOW()
		RETURNING current_page, progress_percent, status
	`
	if err := tx.QueryRowContext(ctx, UpsertReadingProgressQuery, userId, bookId, domain.Reading, domain.Finished).Scan(
		&readingState.CurrentPage,
		&readingState.ProgressPercent,
		&readingState.Status,
	); err != nil {
		rs.l.Error("Error upserting reading progress", "error", err)
		return domain.ReadingState{}, err
	}

	const CloseStaleReadingSessionsQuery = `
		UPDATE reading_sessions
		SET ended_at = $3,
			minutes = GREATEST(
				0,
				LEAST(
					FLOOR(EXTRACT(EPOCH FROM ($3 - started_at)) / 60)::int,
					180
				)
			)
		WHERE user_id = $1
			AND book_id = $2
			AND ended_at IS NULL
	`
	if _, err := tx.ExecContext(ctx, CloseStaleReadingSessionsQuery, userId, bookId, startedAt); err != nil {
		rs.l.Error("Error closing stale reading sessions", "error", err)
		return domain.ReadingState{}, err
	}

	const CreateReadingSessionQuery = `
		INSERT INTO reading_sessions (id, user_id, book_id, started_at, ended_at, minutes)
		VALUES ($1, $2, $3, $4, NULL, 0)
	`
	if _, err := tx.ExecContext(ctx, CreateReadingSessionQuery, sessionId, userId, bookId, startedAt); err != nil {
		rs.l.Error("Error creating reading session", "error", err)
		return domain.ReadingState{}, err
	}
	readingState.SessionId = sessionId

	if err := tx.Commit(); err != nil {
		rs.l.Error("Transaction commit failed", "error", err)
		return domain.ReadingState{}, err
	}
	rs.l.Info("Successfully started reading", "sessionId", readingState.SessionId)
	return readingState, nil
}

func (rs *readingStorage) UpdateProgress(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, currentPage int, progressPercent int, status domain.Status) (domain.ReadingState, error) {
	readingState := domain.ReadingState{BookId: bookId}
	const UpsertProgressQuery = `
		INSERT INTO user_reading_progress (user_id, book_id, status, current_page, progress_percent, updated_at)
		VALUES ($1, $2, $5, $3, $4, NOW())
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET
			current_page = EXCLUDED.current_page,
			progress_percent = EXCLUDED.progress_percent,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING current_page, progress_percent, status
	`
	if err := rs.db.QueryRowContext(ctx, UpsertProgressQuery, userId, bookId, currentPage, progressPercent, status).Scan(
		&readingState.CurrentPage,
		&readingState.ProgressPercent,
		&readingState.Status,
	); err != nil {
		rs.l.Error("Error updating reading progress", "error", err)
		return domain.ReadingState{}, err
	}
	rs.l.Info("Successfully updated reading progress", "bookId", bookId)
	return readingState, nil
}

func (rs *readingStorage) CanRead(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) (bool, error) {
	var canRead bool
	const CanReadQuery = `
		SELECT
			EXISTS (
				SELECT 1
				FROM users_books
				WHERE user_uid = $1 AND book_id = $2
			)
			OR EXISTS (
				SELECT 1
				FROM user_subscriptions
				WHERE user_id = $1
					AND status = 'active'
					AND expired_at > NOW()
			)
	`
	if err := rs.db.QueryRowContext(ctx, CanReadQuery, userId, bookId).Scan(&canRead); err != nil {
		rs.l.Error("Error checking reading access", "error", err)
		return false, err
	}
	rs.l.Info("Successfully checked reading access", "bookId", bookId, "canRead", canRead)
	return canRead, nil
}

func (rs *readingStorage) ActiveSessionBookId(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID) (uuid.UUID, error) {
	var bookId uuid.UUID
	const GetActiveSessionBookQuery = `
		SELECT book_id
		FROM reading_sessions
		WHERE id = $1 AND user_id = $2 AND ended_at IS NULL
	`
	if err := rs.db.QueryRowContext(ctx, GetActiveSessionBookQuery, sessionId, userId).Scan(&bookId); err != nil {
		rs.l.Error("Error getting active reading session book", "error", err)
		return uuid.Nil, err
	}
	rs.l.Info("Successfully got active reading session book", "sessionId", sessionId)
	return bookId, nil
}

func (rs *readingStorage) Finish(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, currentPage int, progressPercent int, status domain.Status, endedAt time.Time) (domain.ReadingState, error) {
	tx, err := rs.db.BeginTx(ctx, nil)
	if err != nil {
		rs.l.Error("Transaction begin failed", "error", err)
		return domain.ReadingState{}, err
	}
	defer tx.Rollback()

	readingState := domain.ReadingState{SessionId: sessionId}
	const FinishReadingSessionQuery = `
		UPDATE reading_sessions
		SET ended_at = $3,
			minutes = GREATEST(0, FLOOR(EXTRACT(EPOCH FROM ($3 - started_at)) / 60)::int)
		WHERE id = $1 AND user_id = $2 AND ended_at IS NULL
		RETURNING book_id
	`
	if err := tx.QueryRowContext(ctx, FinishReadingSessionQuery, sessionId, userId, endedAt).Scan(&readingState.BookId); err != nil {
		rs.l.Error("Error finishing reading session", "error", err)
		return domain.ReadingState{}, err
	}

	const UpsertReadingProgressQuery = `
		INSERT INTO user_reading_progress (user_id, book_id, status, current_page, progress_percent, updated_at)
		VALUES ($1, $2, $5, $3, $4, NOW())
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET
			current_page = EXCLUDED.current_page,
			progress_percent = EXCLUDED.progress_percent,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING current_page, progress_percent, status
	`
	if err := tx.QueryRowContext(ctx, UpsertReadingProgressQuery, userId, readingState.BookId, currentPage, progressPercent, status).Scan(
		&readingState.CurrentPage,
		&readingState.ProgressPercent,
		&readingState.Status,
	); err != nil {
		rs.l.Error("Error upserting reading progress after session finish", "error", err)
		return domain.ReadingState{}, err
	}

	if err := tx.Commit(); err != nil {
		rs.l.Error("Transaction commit failed", "error", err)
		return domain.ReadingState{}, err
	}
	rs.l.Info("Successfully finished reading", "sessionId", sessionId)
	return readingState, nil
}
