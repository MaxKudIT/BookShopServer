package subscription_payments

import (
	"context"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (spserv *subscriptionPaymentsService) Create(ctx context.Context, firebaseId string, subscriptionPayment domain.SubscriptionPayment) (uuid.UUID, error) {
	userId, err := spserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		spserv.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}

	if subscriptionPayment.Currency == "" {
		subscriptionPayment.Currency = "RUB"
	}

	subscriptionPayment.Id = uuid.New()
	subscriptionPayment.UserId = userId
	subscriptionPayment.PaidAt = time.Now()

	if err := spserv.sps.Save(ctx, subscriptionPayment); err != nil {
		spserv.l.Error("Error saving subscription payment", "error", err)
		return subscriptionPayment.Id, err
	}

	spserv.l.Info("Successfully created subscription payment", "id", subscriptionPayment.Id)
	return subscriptionPayment.Id, nil
}

func (spserv *subscriptionPaymentsService) All(ctx context.Context, firebaseId string) ([]domain.SubscriptionPayment, error) {
	userId, err := spserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		spserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	subscriptionPayments, err := spserv.sps.AllByUserId(ctx, userId)
	if err != nil {
		spserv.l.Error("Error getting subscription payments", "error", err)
		return nil, err
	}

	spserv.l.Info("Successfully got subscription payments", "userId", userId)
	return subscriptionPayments, nil
}
