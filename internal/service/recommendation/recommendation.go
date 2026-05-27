package recommendation

import (
	"context"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

const defaultHomeRecommendationLimit = 10

func (rserv *recommendationService) HomeRecommendation(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = defaultHomeRecommendationLimit
	}

	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	books, err := rserv.rs.HomeRecommendation(ctx, userId, limit)
	if err != nil {
		rserv.l.Error("Error getting home recommendations", "error", err)
		return nil, err
	}

	rserv.l.Info("Successfully got home recommendations", "userId", userId)
	return books, nil
}

func (rserv *recommendationService) CartRecommendation(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = defaultHomeRecommendationLimit
	}

	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	books, err := rserv.rs.CartRecommendation(ctx, userId, limit)
	if err != nil {
		rserv.l.Error("Error getting cart recommendations", "error", err)
		return nil, err
	}

	rserv.l.Info("Successfully got cart recommendations", "userId", userId)
	return books, nil
}

func (rserv *recommendationService) RecommesPage(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, []domain.BookPreview, []domain.BookPreview, error) {
	if limit <= 0 {
		limit = defaultHomeRecommendationLimit
	}

	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, nil, nil, err
	}

	forYou, fresh, trend, err := rserv.rs.RecommsPage(ctx, userId, limit)
	if err != nil {
		rserv.l.Error("Error getting recommendations", "error", err)
		return nil, nil, nil, err
	}

	rserv.l.Info("Successfully got recommendations", "userId", userId)
	return forYou, fresh, trend, nil
}

func (rserv *recommendationService) RecommendationByBook(ctx context.Context, firebaseId string, bookId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = defaultHomeRecommendationLimit
	}

	userId, err := rserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		rserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	books, err := rserv.rs.RecommendationByBook(ctx, userId, bookId, limit)
	if err != nil {
		rserv.l.Error("Error getting recommendations by book", "error", err)
		return nil, err
	}

	rserv.l.Info("Successfully got recommendations by book", "userId", userId)
	return books, nil
}
