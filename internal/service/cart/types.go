package cart

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type cartStorage interface {
	Save(ctx context.Context, cart domain.Cart) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type cartService struct {
	cs cartStorage
	us userStorage
	l  *slog.Logger
}

func New(cs cartStorage, us userStorage, l *slog.Logger) *cartService {
	return &cartService{cs: cs, us: us, l: l}
}
