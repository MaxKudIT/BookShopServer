package cart_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type cartItemsService interface {
	IsInCart(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error)
	AllCartItems(ctx context.Context, firebaseId string) ([]domain.CartItemPreview, error)
	Create(ctx context.Context, firebaseId string, cartItem domain.CartItem) (uuid.UUID, error)
	CreateItems(ctx context.Context, firebaseId string, cartItems []domain.CartItem) (uuid.UUID, error)
	AreAllInCart(ctx context.Context, firebaseId string, bookIds []uuid.UUID) (bool, error)
	Delete(ctx context.Context, bookIds []uuid.UUID, firebaseId string) error
	Count(ctx context.Context, firebaseId string) (int, error)
}

type cartItemsHandler struct {
	ciserv cartItemsService
	l      *slog.Logger
}

func New(ciserv cartItemsService, l *slog.Logger) *cartItemsHandler {
	return &cartItemsHandler{ciserv: ciserv, l: l}
}
