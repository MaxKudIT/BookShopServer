package fav_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type favItemsService interface {
	IsInFavs(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error)
	AllFavsItems(ctx context.Context, firebaseId string) ([]domain.FavItemPreview, error)
	Create(ctx context.Context, firebaseId string, favItem domain.FavItem) (uuid.UUID, error)
	Count(ctx context.Context, firebaseId string) (int, error)
	Delete(ctx context.Context, bookIds []uuid.UUID, firebaseId string) error
}

type favItemsHandler struct {
	fiserv favItemsService
	l      *slog.Logger
}

func New(fiserv favItemsService, l *slog.Logger) *favItemsHandler {
	return &favItemsHandler{fiserv: fiserv, l: l}
}
