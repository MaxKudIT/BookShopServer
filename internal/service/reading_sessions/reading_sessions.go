package reading_sessions

import (
	"context"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (rsserv *readingSessionsService) Create(ctx context.Context, readingSession domain.ReadingSession, firebaseId string) (uuid.UUID, error) {
	userId, err := rsserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rsserv.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}
	readingSession.UserId = userId

	if err := rsserv.rss.Save(ctx, readingSession); err != nil {
		rsserv.l.Error("Error saving reading session", "error", err)
		return readingSession.Id, err
	}
	rsserv.l.Info("Successfully created reading session", "id", readingSession.Id)
	return readingSession.Id, nil
}

func (rsserv *readingSessionsService) All(ctx context.Context, firebaseId string) ([]domain.ReadingSession, error) {
	userId, err := rsserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rsserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	readingSessions, err := rsserv.rss.AllByUserId(ctx, userId)
	if err != nil {
		rsserv.l.Error("Error getting reading sessions", "error", err)
		return nil, err
	}
	rsserv.l.Info("Successfully got reading sessions", "userId", userId)
	return readingSessions, nil
}
