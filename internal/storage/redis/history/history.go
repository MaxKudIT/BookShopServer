package history

import (
	"context"
	"fmt"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	bookViewsPrefix    = "book_views"
	readingBooksPrefix = "reading_books"
)

func (hs *historyStorage) SaveBookView(ctx context.Context, userId, bookId uuid.UUID) error {
	return hs.saveBookView(ctx, userId, bookId, time.Now(), 50)
}

func (hs *historyStorage) LastBookViews(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookView, error) {
	if limit <= 0 {
		limit = 50
	}

	values, err := hs.rs.ZRevRangeWithScores(ctx, historyKey(bookViewsPrefix, userId), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	bookViews := make([]domain.BookView, 0, len(values))
	for _, value := range values {
		bookId, err := uuid.Parse(fmt.Sprint(value.Member))
		if err != nil {
			return nil, err
		}

		bookViews = append(bookViews, domain.BookView{
			UserId:   userId,
			BookId:   bookId,
			ViewedAt: time.UnixMilli(int64(value.Score)),
		})
	}

	return bookViews, nil
}

func (hs *historyStorage) WarmBookViews(ctx context.Context, userId uuid.UUID, bookViews []domain.BookView, limit int) error {
	if limit <= 0 {
		limit = 50
	}

	key := historyKey(bookViewsPrefix, userId)
	pipe := hs.rs.TxPipeline()
	pipe.Del(ctx, key)

	for _, bookView := range bookViews {
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(bookView.ViewedAt.UnixMilli()),
			Member: bookView.BookId.String(),
		})
	}

	trimOldRecords(ctx, pipe, key, limit)
	_, err := pipe.Exec(ctx)
	return err
}

func (hs *historyStorage) SaveReadingBook(ctx context.Context, userId, bookId uuid.UUID) error {
	return hs.saveReadingBook(ctx, userId, bookId, time.Now(), 4)
}

func (hs *historyStorage) LastReadingBooks(ctx context.Context, userId uuid.UUID, limit int) ([]domain.LastReadingBook, error) {
	if limit <= 0 {
		limit = 4
	}

	values, err := hs.rs.ZRevRangeWithScores(ctx, historyKey(readingBooksPrefix, userId), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	lastReadingBooks := make([]domain.LastReadingBook, 0, len(values))
	for _, value := range values {
		bookId, err := uuid.Parse(fmt.Sprint(value.Member))
		if err != nil {
			return nil, err
		}

		lastReadingBooks = append(lastReadingBooks, domain.LastReadingBook{
			BookId:        bookId,
			LastStartedAt: time.UnixMilli(int64(value.Score)),
		})
	}

	return lastReadingBooks, nil
}

func (hs *historyStorage) WarmReadingBooks(ctx context.Context, userId uuid.UUID, lastReadingBooks []domain.LastReadingBook, limit int) error {
	if limit <= 0 {
		limit = 4
	}

	key := historyKey(readingBooksPrefix, userId)
	pipe := hs.rs.TxPipeline()
	pipe.Del(ctx, key)

	for _, book := range lastReadingBooks {
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(book.LastStartedAt.UnixMilli()),
			Member: book.BookId.String(),
		})
	}

	trimOldRecords(ctx, pipe, key, limit)
	_, err := pipe.Exec(ctx)
	return err
}

func (hs *historyStorage) saveBookView(ctx context.Context, userId, bookId uuid.UUID, viewedAt time.Time, limit int) error {
	key := historyKey(bookViewsPrefix, userId)

	if err := hs.rs.ZAdd(ctx, key, redis.Z{
		Score:  float64(viewedAt.UnixMilli()),
		Member: bookId.String(),
	}).Err(); err != nil {
		return err
	}

	return hs.rs.ZRemRangeByRank(ctx, key, 0, int64(-(limit + 1))).Err()
}

func (hs *historyStorage) saveReadingBook(ctx context.Context, userId, bookId uuid.UUID, startedAt time.Time, limit int) error {
	key := historyKey(readingBooksPrefix, userId)

	if err := hs.rs.ZAdd(ctx, key, redis.Z{
		Score:  float64(startedAt.UnixMilli()),
		Member: bookId.String(),
	}).Err(); err != nil {
		return err
	}

	return hs.rs.ZRemRangeByRank(ctx, key, 0, int64(-(limit + 1))).Err()
}

func historyKey(prefix string, userId uuid.UUID) string {
	return prefix + ":" + userId.String()
}

func trimOldRecords(ctx context.Context, pipe redis.Pipeliner, key string, limit int) {
	pipe.ZRemRangeByRank(ctx, key, 0, int64(-(limit + 1)))
}
