package user_subscriptions

import (
	"context"
	"errors"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (usserv *userSubscriptionsService) Create(ctx context.Context, firebaseId string, planId uuid.UUID) (uuid.UUID, error) {
	userId, err := usserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		usserv.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}

	if _, isActive, err := usserv.uss.ActiveByUserId(ctx, userId); err != nil {
		usserv.l.Error("Error checking active user subscription", "error", err)
		return uuid.Nil, err
	} else if isActive {
		err := errors.New("user already has active subscription")
		usserv.l.Warn("User already has active subscription", "userId", userId)
		return uuid.Nil, err
	}

	durationDays, err := usserv.uss.PlanDurationDays(ctx, planId)
	if err != nil {
		usserv.l.Error("Error getting subscription plan duration", "error", err)
		return uuid.Nil, err
	}

	startedAt := time.Now()
	userSubscription := domain.UserSubscription{
		Id:        uuid.New(),
		UserId:    userId,
		PlanId:    planId,
		Status:    domain.SubscriptionActive,
		StartedAt: startedAt,
		ExpiresAt: startedAt.AddDate(0, 0, durationDays),
	}

	if err := usserv.uss.Save(ctx, userSubscription); err != nil {
		usserv.l.Error("Error saving user subscription", "error", err)
		return userSubscription.Id, err
	}

	usserv.l.Info("Successfully created user subscription", "id", userSubscription.Id)
	return userSubscription.Id, nil
}

func (usserv *userSubscriptionsService) All(ctx context.Context, firebaseId string) ([]domain.UserSubscription, error) {
	userId, err := usserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		usserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	userSubscriptions, err := usserv.uss.AllByUserId(ctx, userId)
	if err != nil {
		usserv.l.Error("Error getting user subscriptions", "error", err)
		return nil, err
	}

	usserv.l.Info("Successfully got user subscriptions", "userId", userId)
	return userSubscriptions, nil
}

func (usserv *userSubscriptionsService) Status(ctx context.Context, firebaseId string) (domain.UserSubscriptionStatus, error) {
	userId, err := usserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		usserv.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.UserSubscriptionStatus{}, err
	}

	userSubscription, isActive, err := usserv.uss.ActiveByUserId(ctx, userId)
	if err != nil {
		usserv.l.Error("Error getting active user subscription", "error", err)
		return domain.UserSubscriptionStatus{}, err
	}

	if !isActive {
		usserv.l.Info("User subscription is not active", "userId", userId)
		return domain.UserSubscriptionStatus{IsActive: false}, nil
	}

	usserv.l.Info("User subscription is active", "userId", userId)
	return domain.UserSubscriptionStatus{
		IsActive:     true,
		Subscription: &userSubscription,
	}, nil
}
