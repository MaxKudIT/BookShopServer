package cart_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type cartStorage interface {
	CartByUserId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
}

type cartItemStorage interface {
	AllCartItems(ctx context.Context, cartId uuid.UUID) ([]domain.CartItemPreview, error)
	AllCartItemsId(ctx context.Context, cartId uuid.UUID) ([]uuid.UUID, error)
	Save(ctx context.Context, cartItem domain.CartItem) error
	SaveFromFavs(ctx context.Context, cartItems []domain.CartItem) error
	Delete(ctx context.Context, bookIds []uuid.UUID, cartId uuid.UUID) error
	IsInCart(ctx context.Context, cartId uuid.UUID, bookId uuid.UUID) (bool, error)
	AreAllInCart(ctx context.Context, cartId uuid.UUID, bookIds []uuid.UUID) (bool, error)
	Count(ctx context.Context, cartId uuid.UUID) (int, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type cartItemService struct {
	cs  cartStorage
	cis cartItemStorage
	us  userStorage
	l   *slog.Logger
}

func New(cis cartItemStorage, cs cartStorage, us userStorage, l *slog.Logger) *cartItemService {
	return &cartItemService{cis: cis, cs: cs, us: us, l: l}
}
