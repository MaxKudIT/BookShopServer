package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (bs *bookStorage) AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {

	books := make([]domain.BookPreview, 0)
	const AllBooksQuery = `
  SELECT 
    b.id,
    b.title, 
    b.genre, 
    b.price, 
    b.discount, 
    b.image_url,
    b.rate,
    CASE 
        WHEN ub.user_uid IS NOT NULL THEN true
        ELSE false
    END AS is_mine
FROM books b
LEFT JOIN users_books ub ON b.id = ub.book_id 
    AND ub.user_uid = $1;
`
	rows, err := bs.db.QueryContext(ctx, AllBooksQuery, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			bs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			bs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			bs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var currentObject domain.BookPreview
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate, &currentObject.IsMine); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		books = append(books, currentObject)
	}
	fmt.Println(len(books))
	bs.l.Info("Successfully getting books")

	return books, nil
}

func (bs *bookStorage) AllMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {

	myBooks := make([]domain.BookPreview, 0)
	const AllMyBooksQuery = `
    SELECT 
        id, 
        title, 
        genre, 
        price, 
        discount, 
        image_url,
        rate
    FROM books b
    inner join users_books ub on ub.book_id = b.id AND ub.user_uid = $1
`
	rows, err := bs.db.QueryContext(ctx, AllMyBooksQuery, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			bs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			bs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			bs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var currentObject domain.BookPreview
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		myBooks = append(myBooks, currentObject)
	}

	bs.l.Info("Successfully getting my books", "userid", userId)

	return myBooks, nil
}

func (bs *bookStorage) BookById(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) (domain.Book, error) {
	book := domain.Book{}
	const GetBookQuery = `SELECT 
    b.id, 
    b.title, 
    b.author, 
    b.description, 
    b.created_date, 
    b.genre, 
    b.pages_count, 
    b.reading_time, 
    b.price, 
    b.discount, 
    b.about_book, 
    b.quote, 
    b.image_url, 
    b.rate,
    EXISTS(
        SELECT 1 
        FROM users_books ub 
        WHERE ub.book_id = b.id 
          AND ub.user_uid = $1
    ) AS is_mine
FROM books b
WHERE b.id = $2;`
	if err := bs.db.QueryRowContext(ctx, GetBookQuery, userId, bookId).Scan(
		&book.Id,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.CreatedDate,
		&book.Genre,
		&book.PagesCount,
		&book.ReadingTime,
		&book.Price,
		&book.Discount,
		&book.AboutBook,
		&book.Quote,
		&book.ImageUrl,
		&book.Rate,
		&book.IsMine); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			bs.l.Error("book not found", "error", err)
			return domain.Book{}, err
		case errors.Is(err, context.Canceled):
			bs.l.Warn("Query cancelled", "error", err)
			return domain.Book{}, err
		case errors.Is(err, context.DeadlineExceeded):
			bs.l.Warn("Query timed out", "error", err)
			return domain.Book{}, err
		default:
			bs.l.Error("Query failed", "error", err)
			return domain.Book{}, err
		}
	}
	bs.l.Info("Successfully got book", "id", bookId)
	return book, nil
}

func (bs *bookStorage) IsMyBook(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) (bool, error) {
	isMy := false
	const IsMyBookQuery = `
    SELECT	
        EXISTS(
        SELECT 1 
        FROM users_books ub 
        WHERE ub.book_id = $2 
          AND ub.user_uid = $1
    )`
	if err := bs.db.QueryRowContext(ctx, IsMyBookQuery, userId, bookId).Scan(
		&isMy); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			bs.l.Error("book or user not found", "error", err)
			return false, err
		case errors.Is(err, context.Canceled):
			bs.l.Warn("Query cancelled", "error", err)
			return false, err
		case errors.Is(err, context.DeadlineExceeded):
			bs.l.Warn("Query timed out", "error", err)
			return false, err
		default:
			bs.l.Error("Query failed", "error", err)
			return false, err
		}
	}
	bs.l.Info("Successfully got result about purchase", "id", bookId)
	return isMy, nil
}

func (bs *bookStorage) PagesCountById(ctx context.Context, bookId uuid.UUID) (int, error) {
	var pagesCount int
	const GetPagesCountQuery = "SELECT pages_count FROM books WHERE id = $1"
	if err := bs.db.QueryRowContext(ctx, GetPagesCountQuery, bookId).Scan(&pagesCount); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			bs.l.Error("book not found", "error", err)
			return 0, err
		case errors.Is(err, context.Canceled):
			bs.l.Warn("Query cancelled", "error", err)
			return 0, err
		case errors.Is(err, context.DeadlineExceeded):
			bs.l.Warn("Query timed out", "error", err)
			return 0, err
		default:
			bs.l.Error("Query failed", "error", err)
			return 0, err
		}
	}
	bs.l.Info("Successfully got pages count", "bookId", bookId)
	return pagesCount, nil
}
