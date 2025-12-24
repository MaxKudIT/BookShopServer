package user

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type userService interface {
	Create(ctx context.Context, userCr domain.User) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userHandler struct {
	us userService
	l  *slog.Logger
}

func New(us userService, l *slog.Logger) *userHandler {
	return &userHandler{us: us, l: l}
}
