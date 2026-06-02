package orders

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type ordersStorage interface {
	Save(ctx context.Context, order domain.Order) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type ordersService struct {
	os ordersStorage
	us userStorage
	l  *slog.Logger
}

func New(os ordersStorage, us userStorage, l *slog.Logger) *ordersService {
	return &ordersService{os: os, us: us, l: l}
}
