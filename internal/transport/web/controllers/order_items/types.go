package order_items

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type orderItemsService interface {
	Create(ctx context.Context, firebaseId string, orderItem domain.OrderItem) error
	All(ctx context.Context, firebaseId string, orderId uuid.UUID) ([]domain.OrderItem, error)
}

type orderItemsHandler struct {
	oiserv orderItemsService
	l      *slog.Logger
}

func New(oiserv orderItemsService, l *slog.Logger) *orderItemsHandler {
	return &orderItemsHandler{oiserv: oiserv, l: l}
}
