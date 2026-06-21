package reading

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

type fakeBookStorage struct {
	pagesCount int
	err        error
	calls      int
}

func (s *fakeBookStorage) PagesCountById(ctx context.Context, bookId uuid.UUID) (int, error) {
	s.calls++
	if s.err != nil {
		return 0, s.err
	}
	return s.pagesCount, nil
}

type fakeHistoryStorage struct {
	err    error
	calls  int
	userId uuid.UUID
	bookId uuid.UUID
}

func (s *fakeHistoryStorage) SaveReadingBook(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) error {
	s.calls++
	s.userId = userId
	s.bookId = bookId
	return s.err
}

type fakeReadingStorage struct {
	canRead              bool
	canReadErr           error
	startErr             error
	updateErr            error
	activeSessionBookErr error
	finishErr            error
	activeBookId         uuid.UUID

	canReadCalls int
	startCalls   int
	updateCalls  int
	finishCalls  int

	startUserId           uuid.UUID
	startBookId           uuid.UUID
	updateCurrentPage     int
	updateProgressPercent int
	updateStatus          domain.Status
	finishCurrentPage     int
	finishProgressPercent int
	finishStatus          domain.Status
}

func (s *fakeReadingStorage) CanRead(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) (bool, error) {
	s.canReadCalls++
	if s.canReadErr != nil {
		return false, s.canReadErr
	}
	return s.canRead, nil
}

func (s *fakeReadingStorage) Start(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, sessionId uuid.UUID, startedAt time.Time) (domain.ReadingState, error) {
	s.startCalls++
	s.startUserId = userId
	s.startBookId = bookId
	if s.startErr != nil {
		return domain.ReadingState{}, s.startErr
	}
	return domain.ReadingState{SessionId: sessionId, BookId: bookId, Status: domain.Reading, CurrentPage: 1}, nil
}

func (s *fakeReadingStorage) UpdateProgress(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, currentPage int, progressPercent int, status domain.Status) (domain.ReadingState, error) {
	s.updateCalls++
	s.updateCurrentPage = currentPage
	s.updateProgressPercent = progressPercent
	s.updateStatus = status
	if s.updateErr != nil {
		return domain.ReadingState{}, s.updateErr
	}
	return domain.ReadingState{BookId: bookId, Status: status, CurrentPage: currentPage, ProgressPercent: progressPercent}, nil
}

func (s *fakeReadingStorage) ActiveSessionBookId(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID) (uuid.UUID, error) {
	if s.activeSessionBookErr != nil {
		return uuid.Nil, s.activeSessionBookErr
	}
	return s.activeBookId, nil
}

func (s *fakeReadingStorage) Finish(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, currentPage int, progressPercent int, status domain.Status, endedAt time.Time) (domain.ReadingState, error) {
	s.finishCalls++
	s.finishCurrentPage = currentPage
	s.finishProgressPercent = progressPercent
	s.finishStatus = status
	if s.finishErr != nil {
		return domain.ReadingState{}, s.finishErr
	}
	return domain.ReadingState{SessionId: sessionId, BookId: s.activeBookId, Status: status, CurrentPage: currentPage, ProgressPercent: progressPercent}, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCalculateProgress(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		pagesCount  int
		wantPage    int
		wantPercent int
		wantStatus  domain.Status
		wantErr     bool
	}{
		{"first page starts reading", 1, 10, 1, 10, domain.Reading, false},
		{"middle page keeps reading", 5, 10, 5, 50, domain.Reading, false},
		{"last page finishes", 10, 10, 10, 100, domain.Finished, false},
		{"page above total clamps to finish", 15, 10, 10, 100, domain.Finished, false},
		{"single page finishes", 1, 1, 1, 100, domain.Finished, false},
		{"zero current page errors", 0, 10, 0, 0, "", true},
		{"negative current page errors", -1, 10, 0, 0, "", true},
		{"zero pages count errors", 1, 0, 0, 0, "", true},
		{"negative pages count errors", 1, -10, 0, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPercent, gotStatus, err := calculateProgress(tt.currentPage, tt.pagesCount)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPage != tt.wantPage || gotPercent != tt.wantPercent || gotStatus != tt.wantStatus {
				t.Fatalf("got page=%d percent=%d status=%s, want page=%d percent=%d status=%s",
					gotPage, gotPercent, gotStatus, tt.wantPage, tt.wantPercent, tt.wantStatus)
			}
		})
	}
}

func TestReadingServiceStart(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	bookId := uuid.New()

	t.Run("starts reading and stores history", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: true}
		us := &fakeUserStorage{userId: userId}
		hs := &fakeHistoryStorage{}
		svc := New(rs, us, &fakeBookStorage{}, hs, testLogger())

		state, err := svc.Start(ctx, "firebase-id", bookId)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if state.SessionId == uuid.Nil || state.BookId != bookId || state.Status != domain.Reading {
			t.Fatalf("unexpected state: %+v", state)
		}
		if rs.startCalls != 1 || rs.startUserId != userId || rs.startBookId != bookId {
			t.Fatalf("start storage was not called correctly")
		}
		if hs.calls != 1 || hs.userId != userId || hs.bookId != bookId {
			t.Fatalf("history storage was not called correctly")
		}
	})

	t.Run("history failure does not fail start", func(t *testing.T) {
		svc := New(&fakeReadingStorage{canRead: true}, &fakeUserStorage{userId: userId}, &fakeBookStorage{}, &fakeHistoryStorage{err: errors.New("redis down")}, testLogger())
		if _, err := svc.Start(ctx, "firebase-id", bookId); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("denies unavailable book", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: false}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{}, nil, testLogger())
		if _, err := svc.Start(ctx, "firebase-id", bookId); err == nil {
			t.Fatalf("expected error")
		}
		if rs.startCalls != 0 {
			t.Fatalf("start storage should not be called")
		}
	})

	t.Run("user lookup error stops start", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: true}
		svc := New(rs, &fakeUserStorage{err: errors.New("no user")}, &fakeBookStorage{}, nil, testLogger())
		if _, err := svc.Start(ctx, "firebase-id", bookId); err == nil {
			t.Fatalf("expected error")
		}
		if rs.canReadCalls != 0 {
			t.Fatalf("canRead should not be called")
		}
	})

	t.Run("canRead error stops start", func(t *testing.T) {
		rs := &fakeReadingStorage{canReadErr: errors.New("db down")}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{}, nil, testLogger())
		if _, err := svc.Start(ctx, "firebase-id", bookId); err == nil {
			t.Fatalf("expected error")
		}
		if rs.startCalls != 0 {
			t.Fatalf("start storage should not be called")
		}
	})

	t.Run("start storage error returns error", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: true, startErr: errors.New("insert failed")}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{}, nil, testLogger())
		if _, err := svc.Start(ctx, "firebase-id", bookId); err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestReadingServiceUpdateProgress(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	bookId := uuid.New()

	tests := []struct {
		name        string
		page        int
		pagesCount  int
		wantPage    int
		wantPercent int
		wantStatus  domain.Status
		wantErr     bool
	}{
		{"first of ten", 1, 10, 1, 10, domain.Reading, false},
		{"middle progress", 6, 20, 6, 30, domain.Reading, false},
		{"last page finishes", 20, 20, 20, 100, domain.Finished, false},
		{"page over total clamps", 30, 20, 20, 100, domain.Finished, false},
		{"zero page errors", 0, 20, 0, 0, "", true},
		{"zero pages count errors", 1, 0, 0, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &fakeReadingStorage{canRead: true}
			svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{pagesCount: tt.pagesCount}, nil, testLogger())
			_, err := svc.UpdateProgress(ctx, "firebase-id", bookId, tt.page)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if rs.updateCalls != 0 {
					t.Fatalf("update storage should not be called")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rs.updateCurrentPage != tt.wantPage || rs.updateProgressPercent != tt.wantPercent || rs.updateStatus != tt.wantStatus {
				t.Fatalf("got page=%d percent=%d status=%s", rs.updateCurrentPage, rs.updateProgressPercent, rs.updateStatus)
			}
		})
	}

	t.Run("denies progress without access", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: false}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{pagesCount: 10}, nil, testLogger())
		if _, err := svc.UpdateProgress(ctx, "firebase-id", bookId, 2); err == nil {
			t.Fatalf("expected error")
		}
		if rs.updateCalls != 0 {
			t.Fatalf("update storage should not be called")
		}
	})

	t.Run("pages count error stops progress update", func(t *testing.T) {
		rs := &fakeReadingStorage{canRead: true}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{err: errors.New("pages failed")}, nil, testLogger())
		if _, err := svc.UpdateProgress(ctx, "firebase-id", bookId, 2); err == nil {
			t.Fatalf("expected error")
		}
		if rs.updateCalls != 0 {
			t.Fatalf("update storage should not be called")
		}
	})
}

func TestReadingServiceFinish(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	bookId := uuid.New()
	sessionId := uuid.New()

	t.Run("finishes at real pages count with 100 percent", func(t *testing.T) {
		rs := &fakeReadingStorage{activeBookId: bookId}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{pagesCount: 42}, nil, testLogger())
		state, err := svc.Finish(ctx, "firebase-id", sessionId, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if state.Status != domain.Finished || rs.finishCurrentPage != 42 || rs.finishProgressPercent != 100 || rs.finishStatus != domain.Finished {
			t.Fatalf("finish did not force final progress: state=%+v storage=%+v", state, rs)
		}
	})

	t.Run("user lookup error stops finish", func(t *testing.T) {
		rs := &fakeReadingStorage{activeBookId: bookId}
		svc := New(rs, &fakeUserStorage{err: errors.New("no user")}, &fakeBookStorage{pagesCount: 42}, nil, testLogger())
		if _, err := svc.Finish(ctx, "firebase-id", sessionId, 5); err == nil {
			t.Fatalf("expected error")
		}
		if rs.finishCalls != 0 {
			t.Fatalf("finish storage should not be called")
		}
	})

	t.Run("active session lookup error stops finish", func(t *testing.T) {
		rs := &fakeReadingStorage{activeSessionBookErr: errors.New("session missing")}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{pagesCount: 42}, nil, testLogger())
		if _, err := svc.Finish(ctx, "firebase-id", sessionId, 5); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("pages count error stops finish", func(t *testing.T) {
		rs := &fakeReadingStorage{activeBookId: bookId}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{err: errors.New("pages failed")}, nil, testLogger())
		if _, err := svc.Finish(ctx, "firebase-id", sessionId, 5); err == nil {
			t.Fatalf("expected error")
		}
		if rs.finishCalls != 0 {
			t.Fatalf("finish storage should not be called")
		}
	})

	t.Run("finish storage error returns error", func(t *testing.T) {
		rs := &fakeReadingStorage{activeBookId: bookId, finishErr: errors.New("finish failed")}
		svc := New(rs, &fakeUserStorage{userId: userId}, &fakeBookStorage{pagesCount: 42}, nil, testLogger())
		if _, err := svc.Finish(ctx, "firebase-id", sessionId, 5); err == nil {
			t.Fatalf("expected error")
		}
	})
}
