package user_subscriptions

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type fakeUserStorage struct {
	userId uuid.UUID
	err    error
	calls  int
}

func (s *fakeUserStorage) UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error) {
	s.calls++
	if s.err != nil {
		return uuid.Nil, s.err
	}
	return s.userId, nil
}

type fakeUserSubscriptionsStorage struct {
	activeSub       domain.UserSubscription
	active          bool
	activeErr       error
	planDuration    int
	planDurationErr error
	saveErr         error
	all             []domain.UserSubscription
	allErr          error

	saveCalls int
	saved     domain.UserSubscription
	planId    uuid.UUID
}

func (s *fakeUserSubscriptionsStorage) Save(ctx context.Context, userSubscription domain.UserSubscription) error {
	s.saveCalls++
	s.saved = userSubscription
	return s.saveErr
}

func (s *fakeUserSubscriptionsStorage) AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserSubscription, error) {
	if s.allErr != nil {
		return nil, s.allErr
	}
	return s.all, nil
}

func (s *fakeUserSubscriptionsStorage) ActiveByUserId(ctx context.Context, userId uuid.UUID) (domain.UserSubscription, bool, error) {
	if s.activeErr != nil {
		return domain.UserSubscription{}, false, s.activeErr
	}
	return s.activeSub, s.active, nil
}

func (s *fakeUserSubscriptionsStorage) PlanDurationDays(ctx context.Context, planId uuid.UUID) (int, error) {
	s.planId = planId
	if s.planDurationErr != nil {
		return 0, s.planDurationErr
	}
	return s.planDuration, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestUserSubscriptionsServiceCreate(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	planId := uuid.New()

	t.Run("creates active subscription with expiration", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{planDuration: 30}
		svc := New(storage, &fakeUserStorage{userId: userId}, testLogger())
		id, err := svc.Create(ctx, "firebase-id", planId)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id == uuid.Nil || storage.saveCalls != 1 {
			t.Fatalf("subscription was not saved")
		}
		if storage.saved.UserId != userId || storage.saved.PlanId != planId || storage.saved.Status != domain.SubscriptionActive {
			t.Fatalf("unexpected saved subscription: %+v", storage.saved)
		}
		if storage.saved.ExpiresAt.Sub(storage.saved.StartedAt) < 29*24*time.Hour {
			t.Fatalf("expiration is too early: started=%s expires=%s", storage.saved.StartedAt, storage.saved.ExpiresAt)
		}
	})

	t.Run("rejects already active subscription", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{active: true}
		svc := New(storage, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.Create(ctx, "firebase-id", planId); err == nil {
			t.Fatalf("expected error")
		}
		if storage.saveCalls != 0 {
			t.Fatalf("subscription should not be saved")
		}
	})

	t.Run("user lookup error stops create", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{planDuration: 30}
		svc := New(storage, &fakeUserStorage{err: errors.New("no user")}, testLogger())
		if _, err := svc.Create(ctx, "firebase-id", planId); err == nil {
			t.Fatalf("expected error")
		}
		if storage.saveCalls != 0 {
			t.Fatalf("subscription should not be saved")
		}
	})

	t.Run("active lookup error stops create", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{activeErr: errors.New("db down"), planDuration: 30}
		svc := New(storage, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.Create(ctx, "firebase-id", planId); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("plan duration error stops create", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{planDurationErr: errors.New("missing plan")}
		svc := New(storage, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.Create(ctx, "firebase-id", planId); err == nil {
			t.Fatalf("expected error")
		}
		if storage.saveCalls != 0 {
			t.Fatalf("subscription should not be saved")
		}
	})

	t.Run("save error returns generated id and error", func(t *testing.T) {
		storage := &fakeUserSubscriptionsStorage{planDuration: 30, saveErr: errors.New("insert failed")}
		svc := New(storage, &fakeUserStorage{userId: userId}, testLogger())
		id, err := svc.Create(ctx, "firebase-id", planId)
		if err == nil {
			t.Fatalf("expected error")
		}
		if id == uuid.Nil || storage.saveCalls != 1 {
			t.Fatalf("expected generated id and save call")
		}
	})
}

func TestUserSubscriptionsServiceRead(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	activeSub := domain.UserSubscription{Id: uuid.New(), UserId: userId, Status: domain.SubscriptionActive}

	t.Run("status returns inactive when no active subscription", func(t *testing.T) {
		svc := New(&fakeUserSubscriptionsStorage{active: false}, &fakeUserStorage{userId: userId}, testLogger())
		status, err := svc.Status(ctx, "firebase-id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status.IsActive || status.Subscription != nil {
			t.Fatalf("unexpected status: %+v", status)
		}
	})

	t.Run("status returns active subscription", func(t *testing.T) {
		svc := New(&fakeUserSubscriptionsStorage{active: true, activeSub: activeSub}, &fakeUserStorage{userId: userId}, testLogger())
		status, err := svc.Status(ctx, "firebase-id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !status.IsActive || status.Subscription == nil || status.Subscription.Id != activeSub.Id {
			t.Fatalf("unexpected status: %+v", status)
		}
	})

	t.Run("all returns subscriptions", func(t *testing.T) {
		expected := []domain.UserSubscription{{Id: uuid.New()}, {Id: uuid.New()}}
		svc := New(&fakeUserSubscriptionsStorage{all: expected}, &fakeUserStorage{userId: userId}, testLogger())
		got, err := svc.All(ctx, "firebase-id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != len(expected) {
			t.Fatalf("got %d subscriptions, want %d", len(got), len(expected))
		}
	})

	t.Run("all storage error returns error", func(t *testing.T) {
		svc := New(&fakeUserSubscriptionsStorage{allErr: errors.New("db down")}, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.All(ctx, "firebase-id"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("status storage error returns error", func(t *testing.T) {
		svc := New(&fakeUserSubscriptionsStorage{activeErr: errors.New("db down")}, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.Status(ctx, "firebase-id"); err == nil {
			t.Fatalf("expected error")
		}
	})
}
