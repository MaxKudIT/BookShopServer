package order_items

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type orderItemsStorage interface {
	Save(ctx context.Context, userId uuid.UUID, orderItem domain.OrderItem) error
	AllByOrderId(ctx context.Context, userId uuid.UUID, orderId uuid.UUID) ([]domain.OrderItem, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type orderItemsService struct {
	ois orderItemsStorage
	us  userStorage
	l   *slog.Logger
}

func New(ois orderItemsStorage, us userStorage, l *slog.Logger) *orderItemsService {
	return &orderItemsService{ois: ois, us: us, l: l}
}
