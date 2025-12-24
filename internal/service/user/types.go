package user

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type userStorage interface {
	Save(ctx context.Context, userp domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	us userStorage
	l  *slog.Logger
}

func New(us userStorage, l *slog.Logger) *userService {
	return &userService{us: us, l: l}
}
