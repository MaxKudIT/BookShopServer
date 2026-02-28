package fav_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (fiserv *favItemsService) IsInFavs(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error) {

	userId, err := fiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fiserv.l.Error("failed getting id by firebaseId", "err", err)
		return false, err
	}

	favId, err := fiserv.fs.FavByUserId(ctx, userId)
	if err != nil {
		fiserv.l.Info("Failed to get fav id", "id", favId)
		return false, err
	}

	IsInFavs, err := fiserv.fis.IsInFavs(ctx, favId, bookId)
	if err != nil {
		fiserv.l.Error("result about book in the fav failed", "err", err)
		return false, err
	}

	fiserv.l.Info("get a result book in the favs success")
	return IsInFavs, nil
}

func (fiserv *favItemsService) AllFavsItems(ctx context.Context, firebaseId string) ([]domain.FavItemPreview, error) {

	userId, err := fiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fiserv.l.Error("failed getting id by firebaseId", "err", err)
		return nil, err
	}

	favId, err := fiserv.fs.FavByUserId(ctx, userId)
	if err != nil {
		fiserv.l.Info("Failed to get fav id", "id", favId)
		return nil, err
	}

	favItems, err := fiserv.fis.AllFavItems(ctx, favId)

	if err != nil {
		fiserv.l.Error("fav list failed", "err", err)
		return nil, err
	}

	fiserv.l.Info("fav list success")
	return favItems, nil
}

func (fiserv *favItemsService) Create(ctx context.Context, firebaseId string, favItem domain.FavItem) (uuid.UUID, error) {

	userId, err := fiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fiserv.l.Error("failed getting id by firebaseId", "err", err)
		return uuid.Nil, err
	}

	favId, err := fiserv.fs.FavByUserId(ctx, userId)
	if err != nil {
		fiserv.l.Info("Failed to create fav id", "id", favId)
		return uuid.Nil, err
	}

	favItem.FavId = favId

	if err := fiserv.fis.Save(ctx, favItem); err != nil {
		fiserv.l.Error("Error saving fav item", "error", err)
		return favItem.BookId, err
	}
	fiserv.l.Info("Successfully created cart item")
	return favItem.BookId, nil
}

func (fiserv *favItemsService) Count(ctx context.Context, firebaseId string) (int, error) {

	userId, err := fiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fiserv.l.Error("failed getting id by firebaseId", "err", err)
		return 0, err
	}

	favId, err := fiserv.fs.FavByUserId(ctx, userId)
	if err != nil {
		fiserv.l.Info("Failed to get fav id", "id", favId)
		return 0, err
	}

	count, err := fiserv.fis.Count(ctx, favId)
	if err != nil {
		fiserv.l.Error("result about count in the fav failed", "err", err)
		return 0, err
	}

	fiserv.l.Info("get a result count in the fav success")
	return count, nil
}

func (fiserv *favItemsService) Delete(ctx context.Context, bookIds []uuid.UUID, firebaseId string) error {

	userId, err := fiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fiserv.l.Error("failed getting id by firebaseId", "err", err)
		return err
	}

	favId, err := fiserv.fs.FavByUserId(ctx, userId)

	if err != nil {
		fiserv.l.Info("Failed to delete fav id", "id", favId)
		return err
	}

	if err := fiserv.fis.Delete(ctx, bookIds, favId); err != nil {
		fiserv.l.Error("Error deleting fav items", "error", err)
		return err
	}
	fiserv.l.Info("Successfully deleting fav items")
	return nil
}
