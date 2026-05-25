package book_revs

import (
	"context"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (brserv *bookRevsService) Create(ctx context.Context, bookReview domain.BookReview, firebaseId string) (uuid.UUID, error) {
	userId, err := brserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		brserv.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}
	bookReview.UserId = userId

	if err := brserv.brs.Save(ctx, bookReview); err != nil {
		brserv.l.Error("Error saving book review", "error", err)
		return bookReview.Id, err
	}
	brserv.l.Info("Successfully created book review", "id", bookReview.Id)
	return bookReview.Id, nil
}

func (brserv *bookRevsService) All(ctx context.Context, firebaseId string) ([]domain.BookReview, error) {
	userId, err := brserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		brserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	bookReviews, err := brserv.brs.AllByUserId(ctx, userId)
	if err != nil {
		brserv.l.Error("Error getting book reviews", "error", err)
		return nil, err
	}
	brserv.l.Info("Successfully got book reviews", "userId", userId)
	return bookReviews, nil
}
