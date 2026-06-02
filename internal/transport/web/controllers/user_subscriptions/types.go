package user_subscriptions

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type userSubscriptionsService interface {
	Create(ctx context.Context, firebaseId string, planId uuid.UUID) (uuid.UUID, error)
	All(ctx context.Context, firebaseId string) ([]domain.UserSubscription, error)
	Status(ctx context.Context, firebaseId string) (domain.UserSubscriptionStatus, error)
}

type userSubscriptionsHandler struct {
	usserv userSubscriptionsService
	l      *slog.Logger
}

func New(usserv userSubscriptionsService, l *slog.Logger) *userSubscriptionsHandler {
	return &userSubscriptionsHandler{usserv: usserv, l: l}
}
