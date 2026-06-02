package orders

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type ordersService interface {
	Create(ctx context.Context, firebaseId string, order domain.Order) (uuid.UUID, error)
}

type ordersHandler struct {
	oserv ordersService
	l     *slog.Logger
}

func New(oserv ordersService, l *slog.Logger) *ordersHandler {
	return &ordersHandler{oserv: oserv, l: l}
}
