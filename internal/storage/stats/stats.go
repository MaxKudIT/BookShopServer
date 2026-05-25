package stats

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (ss *statsStorage) StatsByUserId(ctx context.Context, userId uuid.UUID) (domain.UserStats, error) {
	var userStats domain.UserStats
	const GetStatsQuery = `
		SELECT
			COALESCE((SELECT SUM(minutes) FROM reading_sessions WHERE user_id = $1), 0),
			COALESCE((SELECT COUNT(*) FROM users_books WHERE user_uid = $1 AND status = 'finished'), 0),
			COALESCE((SELECT AVG(rating) FROM book_revs WHERE user_id = $1), 0)
	`

	if err := ss.db.QueryRowContext(ctx, GetStatsQuery, userId).Scan(
		&userStats.TotalMinutes,
		&userStats.ReadBooks,
		&userStats.AverageRating,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			ss.l.Warn("Query cancelled", "error", err)
			return domain.UserStats{}, err
		case errors.Is(err, context.DeadlineExceeded):
			ss.l.Warn("Query timed out", "error", err)
			return domain.UserStats{}, err
		default:
			ss.l.Error("Query failed", "error", err)
			return domain.UserStats{}, err
		}
	}
	ss.l.Info("Successfully got user stats", "userId", userId)
	return userStats, nil
}
