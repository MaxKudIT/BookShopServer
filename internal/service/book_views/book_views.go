package bookviews

import (
	"context"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (bvserv *bookViewsService) SaveOrUpdate(ctx context.Context, firebaseId string, bookId uuid.UUID) error {
	userId, err := bvserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		bvserv.l.Error("Error getting userId by firebaseId", "error", err)
		return err
	}

	if err := bvserv.bvs.SaveOrUpdate(ctx, userId, bookId); err != nil {
		bvserv.l.Error("Error saving book view", "error", err)
		return err
	}

	if bvserv.hs != nil {
		if err := bvserv.hs.SaveBookView(ctx, userId, bookId); err != nil {
			bvserv.l.Warn("Error saving book view to redis", "error", err)
		}
	}

	bvserv.l.Info("Successfully saved book view", "userId", userId, "bookId", bookId)
	return nil
}

func (bvserv *bookViewsService) LastRecords(ctx context.Context, firebaseId string, limit int) ([]domain.BookViewPreview, error) {
	if limit <= 0 {
		limit = 50
	}

	userId, err := bvserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		bvserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	if bvserv.hs != nil {
		bookViews, err := bvserv.hs.LastBookViews(ctx, userId, limit)
		if err == nil && len(bookViews) > 0 {
			bookPreviews, err := bvserv.bs.BookViewPreviews(ctx, bookViews)
			if err != nil {
				bvserv.l.Error("Error getting book view previews from redis records", "error", err)
				return nil, err
			}

			bvserv.l.Info("Successfully got book view previews from redis", "userId", userId)
			return bookPreviews, nil
		}
		if err != nil {
			bvserv.l.Warn("Error getting book views from redis", "error", err)
		}
	}

	bookViews, err := bvserv.bvs.LastRecords(ctx, userId, limit)
	if err != nil {
		bvserv.l.Error("Error getting book views", "error", err)
		return nil, err
	}

	if bvserv.hs != nil && len(bookViews) > 0 {
		if err := bvserv.hs.WarmBookViews(ctx, userId, bookViews, limit); err != nil {
			bvserv.l.Warn("Error warming book views redis cache", "error", err)
		}
	}

	bookPreviews, err := bvserv.bs.BookViewPreviews(ctx, bookViews)
	if err != nil {
		bvserv.l.Error("Error getting book view previews", "error", err)
		return nil, err
	}

	bvserv.l.Info("Successfully got book view previews", "userId", userId)
	return bookPreviews, nil
}
