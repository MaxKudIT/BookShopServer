package recommendation

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (rs *recommendationStorage) HomeRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	const HomeRecommendationQuery = `
		WITH owned AS (
			SELECT book_id
			FROM users_books
			WHERE user_uid = $1
		),
		seed AS (
			SELECT b.genre, b.author
			FROM books b
			INNER JOIN owned o ON o.book_id = b.id
		),
		genre_counts AS (
			SELECT genre, COUNT(*) AS cnt
			FROM seed
			GROUP BY genre
		),
		author_counts AS (
			SELECT author, COUNT(*) AS cnt
			FROM seed
			GROUP BY author
		)
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate
		FROM books b
		LEFT JOIN genre_counts gc ON gc.genre = b.genre
		LEFT JOIN author_counts ac ON ac.author = b.author
		WHERE NOT EXISTS (
			SELECT 1
			FROM owned o
			WHERE o.book_id = b.id
		)
		ORDER BY
			CASE WHEN gc.cnt IS NOT NULL OR ac.cnt IS NOT NULL THEN 1 ELSE 0 END DESC,
			COALESCE(gc.cnt, 0) * 50 + COALESCE(ac.cnt, 0) * 30 + b.rate * 10 DESC,
			b.rate DESC
		LIMIT $2
	`

	books, err := rs.queryBookPreviews(ctx, HomeRecommendationQuery, userId, limit)
	if err != nil {
		return nil, err
	}

	rs.l.Info("Successfully got home recommendations", "userId", userId)
	return books, nil
}

func (rs *recommendationStorage) CartRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	const CartRecommendationQuery = `
		WITH owned AS (
			SELECT book_id
			FROM users_books
			WHERE user_uid = $1
		),
		cart_books AS (
			SELECT pb.book_id
			FROM cart_items ci
			INNER JOIN physical_books pb ON pb.id = ci.physical_book_id
			INNER JOIN carts c ON c.id = ci.cart_id
			WHERE c.user_id = $1
		),
		seed AS (
			SELECT b.genre, b.author
			FROM books b
			INNER JOIN cart_books cb ON cb.book_id = b.id
		),
		genre_counts AS (
			SELECT genre, COUNT(*) AS cnt
			FROM seed
			GROUP BY genre
		),
		author_counts AS (
			SELECT author, COUNT(*) AS cnt
			FROM seed
			GROUP BY author
		)
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate
		FROM books b
		LEFT JOIN genre_counts gc ON gc.genre = b.genre
		LEFT JOIN author_counts ac ON ac.author = b.author
		WHERE NOT EXISTS (
			SELECT 1
			FROM owned o
			WHERE o.book_id = b.id
		)
		AND NOT EXISTS (
			SELECT 1
			FROM cart_books cb
			WHERE cb.book_id = b.id
		)
		ORDER BY
			CASE WHEN gc.cnt IS NOT NULL OR ac.cnt IS NOT NULL THEN 1 ELSE 0 END DESC,
			COALESCE(gc.cnt, 0) * 50 + COALESCE(ac.cnt, 0) * 30 + b.rate * 10 DESC,
			b.rate DESC
		LIMIT $2
	`

	books, err := rs.queryBookPreviews(ctx, CartRecommendationQuery, userId, limit)
	if err != nil {
		return nil, err
	}

	rs.l.Info("Successfully got cart recommendations", "userId", userId)
	return books, nil
}

func (rs *recommendationStorage) RecommsPage(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, []domain.BookPreview, []domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	forYou, err := rs.HomeRecommendation(ctx, userId, limit)
	if err != nil {
		return nil, nil, nil, err
	}

	fresh, err := rs.FreshRecommendation(ctx, userId, limit)
	if err != nil {
		return nil, nil, nil, err
	}

	trend, err := rs.TrendRecommendation(ctx, userId, limit)
	if err != nil {
		return nil, nil, nil, err
	}

	rs.l.Info("Successfully got recommendations page", "userId", userId)
	return forYou, fresh, trend, nil
}

func (rs *recommendationStorage) FreshRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	const FreshRecommendationQuery = `
		WITH owned AS (
			SELECT book_id
			FROM users_books
			WHERE user_uid = $1
		)
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate
		FROM books b
		WHERE NOT EXISTS (
			SELECT 1
			FROM owned o
			WHERE o.book_id = b.id
		)
		ORDER BY b.created_date DESC, b.rate DESC
		LIMIT $2
	`

	books, err := rs.queryBookPreviews(ctx, FreshRecommendationQuery, userId, limit)
	if err != nil {
		return nil, err
	}

	rs.l.Info("Successfully got fresh recommendations", "userId", userId)
	return books, nil
}

func (rs *recommendationStorage) TrendRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	const TrendRecommendationQuery = `
		WITH owned AS (
			SELECT book_id
			FROM users_books
			WHERE user_uid = $1
		),
		activity AS (
			SELECT book_id, COUNT(*) * 50 AS score
			FROM users_books
			GROUP BY book_id

			UNION ALL

			SELECT pb.book_id, COUNT(*) * 20 AS score
			FROM cart_items ci
			INNER JOIN physical_books pb ON pb.id = ci.physical_book_id
			GROUP BY pb.book_id

			UNION ALL

			SELECT book_id, COUNT(*) * 10 AS score
			FROM reading_sessions
			GROUP BY book_id
		),
		trend_scores AS (
			SELECT book_id, SUM(score) AS score
			FROM activity
			GROUP BY book_id
		)
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate
		FROM books b
		LEFT JOIN trend_scores ts ON ts.book_id = b.id
		WHERE NOT EXISTS (
			SELECT 1
			FROM owned o
			WHERE o.book_id = b.id
		)
		ORDER BY COALESCE(ts.score, 0) DESC, b.rate DESC
		LIMIT $2
	`

	books, err := rs.queryBookPreviews(ctx, TrendRecommendationQuery, userId, limit)
	if err != nil {
		return nil, err
	}

	rs.l.Info("Successfully got trend recommendations", "userId", userId)
	return books, nil
}

func (rs *recommendationStorage) RecommendationByBook(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, limit int) ([]domain.BookPreview, error) {
	if limit <= 0 {
		limit = 10
	}

	const RecommendationByBookQuery = `
		WITH owned AS (
			SELECT book_id
			FROM users_books
			WHERE user_uid = $1
		),
		seed AS (
			SELECT genre, author
			FROM books
			WHERE id = $2
		)
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate
		FROM books b
		CROSS JOIN seed s
		WHERE b.id <> $2
		AND NOT EXISTS (
			SELECT 1
			FROM owned o
			WHERE o.book_id = b.id
		)
		ORDER BY
			CASE WHEN b.genre = s.genre OR b.author = s.author THEN 1 ELSE 0 END DESC,
			CASE WHEN b.genre = s.genre THEN 50 ELSE 0 END +
			CASE WHEN b.author = s.author THEN 30 ELSE 0 END +
			b.rate * 10 DESC,
			b.rate DESC
		LIMIT $3
	`

	books, err := rs.queryBookPreviews(ctx, RecommendationByBookQuery, userId, bookId, limit)
	if err != nil {
		return nil, err
	}

	rs.l.Info("Successfully got book recommendations", "userId", userId, "bookId", bookId)
	return books, nil
}

func (rs *recommendationStorage) queryBookPreviews(ctx context.Context, query string, args ...any) ([]domain.BookPreview, error) {
	books := make([]domain.BookPreview, 0)

	rows, err := rs.db.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			rs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			rs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			rs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.BookPreview
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.Title,
			&currentObject.Author,
			&currentObject.Genre,
			&currentObject.Price,
			&currentObject.Discount,
			&currentObject.ImageUrl,
			&currentObject.Rate,
		); err != nil {
			rs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		books = append(books, currentObject)
	}

	if err := rows.Err(); err != nil {
		rs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	return books, nil
}
