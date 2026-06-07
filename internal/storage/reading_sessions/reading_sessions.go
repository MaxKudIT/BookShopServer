package reading_sessions

import (
	"context"
	"errors"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (rss *readingSessionsStorage) Save(ctx context.Context, readingSession domain.ReadingSession) error {
	const CreateReadingSessionQuery = `
		INSERT INTO reading_sessions (id, user_id, book_id, started_at, ended_at, minutes)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if _, err := rss.db.ExecContext(
		ctx,
		CreateReadingSessionQuery,
		readingSession.Id,
		readingSession.UserId,
		readingSession.BookId,
		readingSession.StartedAt,
		readingSession.EndedAt,
		readingSession.Minutes,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			rss.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			rss.l.Warn("Query timed out", "error", err)
			return err
		default:
			rss.l.Error("Query failed", "error", err)
			return err
		}
	}
	rss.l.Info("Successfully saved reading session", "id", readingSession.Id)
	return nil
}

func (rss *readingSessionsStorage) AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.ReadingSession, error) {
	readingSessions := make([]domain.ReadingSession, 0)
	const AllReadingSessionsQuery = `
		SELECT id, user_id, book_id, started_at, ended_at, minutes
		FROM reading_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
	`

	rows, err := rss.db.QueryContext(ctx, AllReadingSessionsQuery, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			rss.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			rss.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			rss.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.ReadingSession
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.UserId,
			&currentObject.BookId,
			&currentObject.StartedAt,
			&currentObject.EndedAt,
			&currentObject.Minutes,
		); err != nil {
			rss.l.Error("Scan failed", "error", err)
			return nil, err
		}
		readingSessions = append(readingSessions, currentObject)
	}

	if err := rows.Err(); err != nil {
		rss.l.Error("Rows failed", "error", err)
		return nil, err
	}

	rss.l.Info("Successfully got reading sessions", "userId", userId)
	return readingSessions, nil
}

func (rss *readingSessionsStorage) Close(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, endedAt time.Time) (domain.ReadingSession, error) {
	var readingSession domain.ReadingSession
	const CloseReadingSessionQuery = `
		UPDATE reading_sessions
		SET ended_at = $3,
			minutes = GREATEST(0, FLOOR(EXTRACT(EPOCH FROM ($3 - started_at)) / 60)::int)
		WHERE id = $1
			AND user_id = $2
			AND ended_at IS NULL
		RETURNING id, user_id, book_id, started_at, ended_at, minutes
	`

	if err := rss.db.QueryRowContext(ctx, CloseReadingSessionQuery, sessionId, userId, endedAt).Scan(
		&readingSession.Id,
		&readingSession.UserId,
		&readingSession.BookId,
		&readingSession.StartedAt,
		&readingSession.EndedAt,
		&readingSession.Minutes,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			rss.l.Warn("Query cancelled", "error", err)
			return domain.ReadingSession{}, err
		case errors.Is(err, context.DeadlineExceeded):
			rss.l.Warn("Query timed out", "error", err)
			return domain.ReadingSession{}, err
		default:
			rss.l.Error("Query failed", "error", err)
			return domain.ReadingSession{}, err
		}
	}

	rss.l.Info("Successfully closed reading session", "id", sessionId)
	return readingSession, nil
}

func (rss *readingSessionsStorage) LastReadingBookRecords(ctx context.Context, userId uuid.UUID, limit int) ([]domain.LastReadingBook, error) {
	if limit <= 0 {
		limit = 4
	}

	lastReadingBooks := make([]domain.LastReadingBook, 0, limit)
	const LastReadingBookRecordsQuery = `
		SELECT book_id, MAX(started_at) AS last_started_at
		FROM reading_sessions
		WHERE user_id = $1
		GROUP BY book_id
		ORDER BY last_started_at DESC
		LIMIT $2
	`

	rows, err := rss.db.QueryContext(ctx, LastReadingBookRecordsQuery, userId, limit)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			rss.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			rss.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			rss.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.LastReadingBook
		if err := rows.Scan(
			&currentObject.BookId,
			&currentObject.LastStartedAt,
		); err != nil {
			rss.l.Error("Scan failed", "error", err)
			return nil, err
		}
		lastReadingBooks = append(lastReadingBooks, currentObject)
	}

	if err := rows.Err(); err != nil {
		rss.l.Error("Rows failed", "error", err)
		return nil, err
	}

	rss.l.Info("Successfully got last reading book records", "userId", userId)
	return lastReadingBooks, nil
}
