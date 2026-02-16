package cart

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type cartService interface {
	Create(ctx context.Context, cart domain.Cart, firebaseId string) (uuid.UUID, error)
}

type cartHandler struct {
	cs cartService
	l  *slog.Logger
}

func New(cs cartService, l *slog.Logger) *cartHandler {
	return &cartHandler{cs: cs, l: l}
}
