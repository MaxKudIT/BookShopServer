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

func (rsserv *readingSessionsService) LastReadingBooks(ctx context.Context, firebaseId string, limit int) ([]domain.ReadingBookPreview, error) {
	if limit <= 0 {
		limit = 4
	}

	userId, err := rsserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rsserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	if rsserv.hs != nil {
		lastReadingBooks, err := rsserv.hs.LastReadingBooks(ctx, userId, limit)
		if err == nil && len(lastReadingBooks) > 0 {
			books, err := rsserv.bs.ReadingBookPreviews(ctx, lastReadingBooks)
			if err != nil {
				rsserv.l.Error("Error getting reading book previews from redis records", "error", err)
				return nil, err
			}

			rsserv.l.Info("Successfully got last reading books from redis", "userId", userId)
			return books, nil
		}
		if err != nil {
			rsserv.l.Warn("Error getting last reading books from redis", "error", err)
		}
	}

	lastReadingBooks, err := rsserv.rss.LastReadingBookRecords(ctx, userId, limit)
	if err != nil {
		rsserv.l.Error("Error getting last reading book records", "error", err)
		return nil, err
	}

	if rsserv.hs != nil && len(lastReadingBooks) > 0 {
		if err := rsserv.hs.WarmReadingBooks(ctx, userId, lastReadingBooks, limit); err != nil {
			rsserv.l.Warn("Error warming last reading books redis cache", "error", err)
		}
	}

	books, err := rsserv.bs.ReadingBookPreviews(ctx, lastReadingBooks)
	if err != nil {
		rsserv.l.Error("Error getting reading book previews", "error", err)
		return nil, err
	}

	rsserv.l.Info("Successfully got last reading books", "userId", userId)
	return books, nil
}
