package recommendation

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

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

type fakeRecommendationStorage struct {
	homeBooks []domain.BookPreview
	cartBooks []domain.BookPreview
	forYou    []domain.BookPreview
	fresh     []domain.BookPreview
	trend     []domain.BookPreview
	byBook    []domain.BookPreview

	homeErr   error
	cartErr   error
	pageErr   error
	byBookErr error

	homeLimit   int
	cartLimit   int
	pageLimit   int
	byBookLimit int
	userId      uuid.UUID
	bookId      uuid.UUID
}

func (s *fakeRecommendationStorage) HomeRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	s.userId = userId
	s.homeLimit = limit
	if s.homeErr != nil {
		return nil, s.homeErr
	}
	return s.homeBooks, nil
}

func (s *fakeRecommendationStorage) CartRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	s.userId = userId
	s.cartLimit = limit
	if s.cartErr != nil {
		return nil, s.cartErr
	}
	return s.cartBooks, nil
}

func (s *fakeRecommendationStorage) RecommsPage(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, []domain.BookPreview, []domain.BookPreview, error) {
	s.userId = userId
	s.pageLimit = limit
	if s.pageErr != nil {
		return nil, nil, nil, s.pageErr
	}
	return s.forYou, s.fresh, s.trend, nil
}

func (s *fakeRecommendationStorage) RecommendationByBook(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	s.userId = userId
	s.bookId = bookId
	s.byBookLimit = limit
	if s.byBookErr != nil {
		return nil, s.byBookErr
	}
	return s.byBook, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func sampleBooks(count int) []domain.BookPreview {
	books := make([]domain.BookPreview, 0, count)
	for i := 0; i < count; i++ {
		books = append(books, domain.BookPreview{Id: uuid.New(), Title: "Book"})
	}
	return books
}

func TestRecommendationServiceHome(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	t.Run("uses explicit limit and user id", func(t *testing.T) {
		rs := &fakeRecommendationStorage{homeBooks: sampleBooks(2)}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		books, err := svc.HomeRecommendation(ctx, "firebase-id", 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 2 || rs.homeLimit != 4 || rs.userId != userId {
			t.Fatalf("unexpected home recommendation call")
		}
	})

	t.Run("uses default limit for invalid limit", func(t *testing.T) {
		rs := &fakeRecommendationStorage{}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.HomeRecommendation(ctx, "firebase-id", 0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rs.homeLimit != defaultHomeRecommendationLimit {
			t.Fatalf("got limit %d, want %d", rs.homeLimit, defaultHomeRecommendationLimit)
		}
	})

	t.Run("user lookup error returns error", func(t *testing.T) {
		rs := &fakeRecommendationStorage{}
		svc := New(rs, &fakeUserStorage{err: errors.New("no user")}, testLogger())
		if _, err := svc.HomeRecommendation(ctx, "firebase-id", 4); err == nil {
			t.Fatalf("expected error")
		}
		if rs.homeLimit != 0 {
			t.Fatalf("storage should not be called")
		}
	})

	t.Run("storage error returns error", func(t *testing.T) {
		svc := New(&fakeRecommendationStorage{homeErr: errors.New("db down")}, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.HomeRecommendation(ctx, "firebase-id", 4); err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestRecommendationServiceOtherFlows(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	bookId := uuid.New()

	t.Run("cart recommendations pass limit", func(t *testing.T) {
		rs := &fakeRecommendationStorage{cartBooks: sampleBooks(3)}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		books, err := svc.CartRecommendation(ctx, "firebase-id", 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 3 || rs.cartLimit != 7 || rs.userId != userId {
			t.Fatalf("unexpected cart recommendation call")
		}
	})

	t.Run("cart recommendations use default limit", func(t *testing.T) {
		rs := &fakeRecommendationStorage{}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.CartRecommendation(ctx, "firebase-id", -1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rs.cartLimit != defaultHomeRecommendationLimit {
			t.Fatalf("got limit %d", rs.cartLimit)
		}
	})

	t.Run("recommendations page returns all sections", func(t *testing.T) {
		rs := &fakeRecommendationStorage{
			forYou: sampleBooks(1),
			fresh:  sampleBooks(2),
			trend:  sampleBooks(3),
		}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		forYou, fresh, trend, err := svc.RecommesPage(ctx, "firebase-id", 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(forYou) != 1 || len(fresh) != 2 || len(trend) != 3 || rs.pageLimit != 5 {
			t.Fatalf("unexpected recommendations page result")
		}
	})

	t.Run("recommendation by book passes book id", func(t *testing.T) {
		rs := &fakeRecommendationStorage{byBook: sampleBooks(2)}
		svc := New(rs, &fakeUserStorage{userId: userId}, testLogger())
		books, err := svc.RecommendationByBook(ctx, "firebase-id", bookId, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 2 || rs.bookId != bookId || rs.byBookLimit != 6 {
			t.Fatalf("unexpected by-book recommendation call")
		}
	})

	t.Run("recommendations page storage error returns error", func(t *testing.T) {
		svc := New(&fakeRecommendationStorage{pageErr: errors.New("db down")}, &fakeUserStorage{userId: userId}, testLogger())
		if _, _, _, err := svc.RecommesPage(ctx, "firebase-id", 5); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("recommendation by book storage error returns error", func(t *testing.T) {
		svc := New(&fakeRecommendationStorage{byBookErr: errors.New("db down")}, &fakeUserStorage{userId: userId}, testLogger())
		if _, err := svc.RecommendationByBook(ctx, "firebase-id", bookId, 5); err == nil {
			t.Fatalf("expected error")
		}
	})
}
