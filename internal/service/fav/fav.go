package fav

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (fserv *favService) Create(ctx context.Context, fav domain.Fav, firebaseId string) (uuid.UUID, error) {

	userId, err := fserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		fserv.l.Error("fav saving failed", "err", err)
		return uuid.Nil, err
	}
	fav.UserId = userId

	if err := fserv.fs.Save(ctx, fav); err != nil {
		fserv.l.Error("Error saving fav", "error", err)
		return fav.Id, err
	}
	fserv.l.Info("Successfully created fav")
	return fav.Id, nil
}
