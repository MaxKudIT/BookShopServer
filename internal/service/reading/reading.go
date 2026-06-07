package reading

import (
	"context"
	"errors"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (rserv *readingService) Start(ctx context.Context, firebaseId string, bookId uuid.UUID) (domain.ReadingState, error) {
	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.ReadingState{}, err
	}

	canRead, err := rserv.rs.CanRead(ctx, userId, bookId)
	if err != nil {
		rserv.l.Error("Error checking reading access", "error", err)
		return domain.ReadingState{}, err
	}
	if !canRead {
		rserv.l.Warn("Reading access denied", "userId", userId, "bookId", bookId)
		return domain.ReadingState{}, errors.New("book is not available for reading")
	}

	readingState, err := rserv.rs.Start(ctx, userId, bookId, uuid.New(), time.Now())
	if err != nil {
		rserv.l.Error("Error starting reading", "error", err)
		return domain.ReadingState{}, err
	}

	if rserv.hs != nil {
		if err := rserv.hs.SaveReadingBook(ctx, userId, bookId); err != nil {
			rserv.l.Warn("Error saving last reading book to redis", "error", err)
		}
	}

	rserv.l.Info("Successfully started reading", "sessionId", readingState.SessionId)
	return readingState, nil
}

func (rserv *readingService) UpdateProgress(ctx context.Context, firebaseId string, bookId uuid.UUID, currentPage int) (domain.ReadingState, error) {
	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.ReadingState{}, err
	}

	canRead, err := rserv.rs.CanRead(ctx, userId, bookId)
	if err != nil {
		rserv.l.Error("Error checking reading access", "error", err)
		return domain.ReadingState{}, err
	}
	if !canRead {
		rserv.l.Warn("Reading access denied", "userId", userId, "bookId", bookId)
		return domain.ReadingState{}, errors.New("book is not available for reading")
	}

	pagesCount, err := rserv.bs.PagesCountById(ctx, bookId)
	if err != nil {
		rserv.l.Error("Error getting pages count", "error", err)
		return domain.ReadingState{}, err
	}

	currentPage, progressPercent, status, err := calculateProgress(currentPage, pagesCount)
	if err != nil {
		rserv.l.Error("Error calculating progress", "error", err)
		return domain.ReadingState{}, err
	}

	readingState, err := rserv.rs.UpdateProgress(ctx, userId, bookId, currentPage, progressPercent, status)
	if err != nil {
		rserv.l.Error("Error updating reading progress", "error", err)
		return domain.ReadingState{}, err
	}
	rserv.l.Info("Successfully updated reading progress", "bookId", bookId)
	return readingState, nil
}

func (rserv *readingService) Finish(ctx context.Context, firebaseId string, sessionId uuid.UUID, currentPage int) (domain.ReadingState, error) {
	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.ReadingState{}, err
	}

	bookId, err := rserv.rs.ActiveSessionBookId(ctx, userId, sessionId)
	if err != nil {
		rserv.l.Error("Error getting active reading session book", "error", err)
		return domain.ReadingState{}, err
	}

	pagesCount, err := rserv.bs.PagesCountById(ctx, bookId)
	if err != nil {
		rserv.l.Error("Error getting pages count", "error", err)
		return domain.ReadingState{}, err
	}

	currentPage, progressPercent, status, err := calculateProgress(currentPage, pagesCount)
	if err != nil {
		rserv.l.Error("Error calculating progress", "error", err)
		return domain.ReadingState{}, err
	}

	readingState, err := rserv.rs.Finish(ctx, userId, sessionId, currentPage, progressPercent, status, time.Now())
	if err != nil {
		rserv.l.Error("Error finishing reading", "error", err)
		return domain.ReadingState{}, err
	}
	rserv.l.Info("Successfully finished reading", "sessionId", sessionId)
	return readingState, nil
}

func calculateProgress(currentPage int, pagesCount int) (int, int, domain.Status, error) {
	if pagesCount <= 0 {
		return 0, 0, "", errors.New("pages count must be greater than zero")
	}
	if currentPage < 1 {
		return 0, 0, "", errors.New("current page must be greater than zero")
	}
	if currentPage > pagesCount {
		currentPage = pagesCount
	}

	progressPercent := currentPage * 100 / pagesCount
	status := domain.Reading
	if progressPercent >= 100 {
		status = domain.Finished
	}

	return currentPage, progressPercent, status, nil
}
