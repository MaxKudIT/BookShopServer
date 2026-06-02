package user_subscriptions

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type userSubscriptionsStorage interface {
	Save(ctx context.Context, userSubscription domain.UserSubscription) error
	AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserSubscription, error)
	ActiveByUserId(ctx context.Context, userId uuid.UUID) (domain.UserSubscription, bool, error)
	PlanDurationDays(ctx context.Context, planId uuid.UUID) (int, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type userSubscriptionsService struct {
	uss userSubscriptionsStorage
	us  userStorage
	l   *slog.Logger
}

func New(uss userSubscriptionsStorage, us userStorage, l *slog.Logger) *userSubscriptionsService {
	return &userSubscriptionsService{uss: uss, us: us, l: l}
}
