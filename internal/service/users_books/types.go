package users_books

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type ubStorage interface {
	Buy(ctx context.Context, userId uuid.UUID, bookIds []uuid.UUID) error
}

type ubService struct {
	ubs ubStorage
	us  userStorage
	l   *slog.Logger
}

func New(ubs ubStorage, us userStorage, l *slog.Logger) *ubService {
	return &ubService{ubs: ubs, us: us, l: l}
}
