package reading

import (
	"context"
	"database/sql"
	"errors"
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
	const UpdateUserBookQuery = `
		UPDATE users_books
		SET status = CASE WHEN status = $3 THEN status ELSE $4 END
		WHERE user_uid = $1 AND book_id = $2
		RETURNING COALESCE(current_page, 0), COALESCE(progress_percent, 0), status
	`
	if err := tx.QueryRowContext(ctx, UpdateUserBookQuery, userId, bookId, domain.Finished, domain.Reading).Scan(
		&readingState.CurrentPage,
		&readingState.ProgressPercent,
		&readingState.Status,
	); err != nil {
		rs.l.Error("Error updating user book status", "error", err)
		return domain.ReadingState{}, err
	}

	const GetActiveSessionQuery = `
		SELECT id
		FROM reading_sessions
		WHERE user_id = $1 AND book_id = $2 AND ended_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1
	`
	if err := tx.QueryRowContext(ctx, GetActiveSessionQuery, userId, bookId).Scan(&readingState.SessionId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			rs.l.Error("Error getting active reading session", "error", err)
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
	}

	if err := tx.Commit(); err != nil {
		rs.l.Error("Transaction commit failed", "error", err)
		return domain.ReadingState{}, err
	}
	rs.l.Info("Successfully started reading", "sessionId", readingState.SessionId)
	return readingState, nil
}

func (rs *readingStorage) UpdateProgress(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, currentPage int, progressPercent int, status domain.Status) (domain.ReadingState, error) {
	readingState := domain.ReadingState{BookId: bookId}
	const UpdateProgressQuery = `
		UPDATE users_books
		SET current_page = $3,
			progress_percent = $4,
			status = $5
		WHERE user_uid = $1 AND book_id = $2
		RETURNING current_page, progress_percent, status
	`
	if err := rs.db.QueryRowContext(ctx, UpdateProgressQuery, userId, bookId, currentPage, progressPercent, status).Scan(
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

	const UpdateUserBookQuery = `
		UPDATE users_books
		SET current_page = $3,
			progress_percent = $4,
			status = $5
		WHERE user_uid = $1 AND book_id = $2
		RETURNING current_page, progress_percent, status
	`
	if err := tx.QueryRowContext(ctx, UpdateUserBookQuery, userId, readingState.BookId, currentPage, progressPercent, status).Scan(
		&readingState.CurrentPage,
		&readingState.ProgressPercent,
		&readingState.Status,
	); err != nil {
		rs.l.Error("Error updating user book after session finish", "error", err)
		return domain.ReadingState{}, err
	}

	if err := tx.Commit(); err != nil {
		rs.l.Error("Transaction commit failed", "error", err)
		return domain.ReadingState{}, err
	}
	rs.l.Info("Successfully finished reading", "sessionId", sessionId)
	return readingState, nil
}
