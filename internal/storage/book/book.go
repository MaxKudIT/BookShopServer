package book

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (bs *bookStorage) AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {

	books := make([]domain.BookPreview, 0)
	const AllBooksQuery = `
    SELECT 
        id, 
        title, 
        genre, 
        price, 
        discount, 
        image_url,
        case 
            when user_id = $1 then true
            else false
        end as is_mine
    FROM books
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
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.IsMine); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		books = append(books, currentObject)
	}

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
        image_url
    FROM books 
    where user_id = $1
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
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.IsMine); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		myBooks = append(myBooks, currentObject)
	}

	bs.l.Info("Successfully getting my books", "userid", userId)

	return myBooks, nil
}

func (cs *bookStorage) BookById(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	book := domain.Book{}
	const GetBookQuery = "SELECT id, title, author, description, created_date, genre, pages_count, reading_time, price, discount, about_book, quote, image_url, rate, case when user_id = $1 then true else false end as is_mine FROM books WHERE id = $2"
	if err := cs.db.QueryRowContext(ctx, GetBookQuery, id).Scan(
		&book.Id,
		book.Title,
		book.Author,
		book.Description,
		book.CreatedDate,
		book.Genre,
		book.PagesCount,
		book.ReadingTime,
		book.Price,
		book.Discount,
		book.AboutBook,
		book.Quote,
		book.ImageUrl,
		book.Rate,
		book.IsMine); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			cs.l.Error("book not found", "error", err)
			return domain.Book{}, err
		case errors.Is(err, context.Canceled):
			cs.l.Warn("Query cancelled", "error", err)
			return domain.Book{}, err
		case errors.Is(err, context.DeadlineExceeded):
			cs.l.Warn("Query timed out", "error", err)
			return domain.Book{}, err
		default:
			cs.l.Error("Query failed", "error", err)
			return domain.Book{}, err
		}
	}
	cs.l.Info("Successfully got book", "id", id)
	return book, nil
}
