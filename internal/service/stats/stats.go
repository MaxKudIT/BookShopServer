package stats

import (
	"context"

	"github.com/bookshop/internal/domain"
)

func (sserv *statsService) UserStats(ctx context.Context, firebaseId string) (domain.UserStats, error) {
	userId, err := sserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		sserv.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.UserStats{}, err
	}

	userStats, err := sserv.ss.StatsByUserId(ctx, userId)
	if err != nil {
		sserv.l.Error("Error getting user stats", "error", err)
		return domain.UserStats{}, err
	}
	sserv.l.Info("Successfully got user stats", "userId", userId)
	return userStats, nil
}
