package stats

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type statsStorage interface {
	StatsByUserId(ctx context.Context, userId uuid.UUID) (domain.UserStats, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type statsService struct {
	ss statsStorage
	us userStorage
	l  *slog.Logger
}

func New(ss statsStorage, us userStorage, l *slog.Logger) *statsService {
	return &statsService{ss: ss, us: us, l: l}
}
