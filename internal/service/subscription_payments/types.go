package subscription_payments

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type subscriptionPaymentsStorage interface {
	Save(ctx context.Context, subscriptionPayment domain.SubscriptionPayment) error
	AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.SubscriptionPayment, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type subscriptionPaymentsService struct {
	sps subscriptionPaymentsStorage
	us  userStorage
	l   *slog.Logger
}

func New(sps subscriptionPaymentsStorage, us userStorage, l *slog.Logger) *subscriptionPaymentsService {
	return &subscriptionPaymentsService{sps: sps, us: us, l: l}
}
