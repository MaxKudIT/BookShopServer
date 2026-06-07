package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (bs *bookStorage) AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {

	books := make([]domain.BookPreview, 0)
	const AllBooksQuery = `
  SELECT 
    b.id,
    b.title, 
	b.author,
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
		var isMine bool
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Author, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate, &isMine); err != nil {
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
		author,
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
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Author, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		myBooks = append(myBooks, currentObject)
	}

	bs.l.Info("Successfully getting my books", "userid", userId)

	return myBooks, nil
}

func (bs *bookStorage) AllNotMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {

	myBooks := make([]domain.BookPreview, 0)
	const AllNotMyBooksQuery = `
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
    FROM users_books ub
    WHERE ub.book_id = b.id
      AND ub.user_uid = $1
);
`
	rows, err := bs.db.QueryContext(ctx, AllNotMyBooksQuery, userId)
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
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Author, &currentObject.Genre, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		myBooks = append(myBooks, currentObject)
	}

	bs.l.Info("Successfully getting not my books", "userid", userId)

	return myBooks, nil
}

func (bs *bookStorage) Search(ctx context.Context, userId uuid.UUID, filter domain.BookSearchFilter) ([]domain.BookPreview, error) {
	books := make([]domain.BookPreview, 0)

	args := []any{userId}
	conditions := []string{"1 = 1"}
	searchPlaceholder := ""

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		searchPlaceholder = fmt.Sprintf("$%d", len(args))
		conditions = append(conditions, fmt.Sprintf(`(
			b.title ILIKE %[1]s
			OR b.author ILIKE %[1]s
			OR b.genre::text ILIKE %[1]s
			OR b.description ILIKE %[1]s
			OR b.about_book ILIKE %[1]s
		)`, searchPlaceholder))
	}

	if filter.Genre != "" {
		args = append(args, filter.Genre)
		conditions = append(conditions, fmt.Sprintf("LOWER(b.genre::text) = LOWER($%d)", len(args)))
	}

	if filter.MinPrice != nil {
		args = append(args, *filter.MinPrice)
		conditions = append(conditions, fmt.Sprintf("b.price >= $%d", len(args)))
	}

	if filter.MaxPrice != nil {
		args = append(args, *filter.MaxPrice)
		conditions = append(conditions, fmt.Sprintf("b.price <= $%d", len(args)))
	}

	if filter.MinRate != nil {
		args = append(args, *filter.MinRate)
		conditions = append(conditions, fmt.Sprintf("b.rate >= $%d", len(args)))
	}

	orderBy := searchOrderBy(filter.Sort, searchPlaceholder)

	args = append(args, filter.Limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))
	args = append(args, filter.Offset)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args))

	SearchBooksQuery := fmt.Sprintf(`
		SELECT
			b.id,
			b.title,
			b.author,
			b.genre,
			b.price,
			b.discount,
			b.image_url,
			b.rate,
			EXISTS(
				SELECT 1
				FROM users_books ub
				WHERE ub.book_id = b.id
				  AND ub.user_uid = $1
			) AS is_mine
		FROM books b
		WHERE %s
		ORDER BY %s
		LIMIT %s OFFSET %s
	`, strings.Join(conditions, " AND "), orderBy, limitPlaceholder, offsetPlaceholder)

	rows, err := bs.db.QueryContext(ctx, SearchBooksQuery, args...)
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
		var isMine bool
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.Title,
			&currentObject.Author,
			&currentObject.Genre,
			&currentObject.Price,
			&currentObject.Discount,
			&currentObject.ImageUrl,
			&currentObject.Rate,
			&isMine,
		); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		books = append(books, currentObject)
	}

	if err := rows.Err(); err != nil {
		bs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	bs.l.Info("Successfully searched books", "count", len(books))
	return books, nil
}

func searchOrderBy(sort string, searchPlaceholder string) string {
	switch sort {
	case "price_asc":
		return "b.price ASC, b.rate DESC, b.title ASC"
	case "price_desc":
		return "b.price DESC, b.rate DESC, b.title ASC"
	case "new":
		return "b.created_date DESC, b.rate DESC, b.title ASC"
	case "title":
		return "b.title ASC, b.rate DESC"
	case "rate_desc":
		return "b.rate DESC, b.title ASC"
	default:
		if searchPlaceholder != "" {
			return fmt.Sprintf("CASE WHEN b.title ILIKE %[1]s THEN 0 WHEN b.author ILIKE %[1]s THEN 1 ELSE 2 END, b.rate DESC, b.title ASC", searchPlaceholder)
		}
		return "b.rate DESC, b.title ASC"
	}
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

func (bs *bookStorage) BookViewPreviews(ctx context.Context, bookViews []domain.BookView) ([]domain.BookViewPreview, error) {
	if len(bookViews) == 0 {
		return make([]domain.BookViewPreview, 0), nil
	}

	bookIds := make([]uuid.UUID, 0, len(bookViews))
	viewedAtByBookId := make(map[uuid.UUID]time.Time, len(bookViews))
	for _, bookView := range bookViews {
		bookIds = append(bookIds, bookView.BookId)
		viewedAtByBookId[bookView.BookId] = bookView.ViewedAt
	}

	previews := make([]domain.BookViewPreview, 0, len(bookViews))
	const BookViewPreviewsQuery = `
		SELECT
			id,
			image_url,
			title,
			author,
			created_date,
			genre,
			rate,
			pages_count
		FROM books
		WHERE id = ANY($1)
		ORDER BY array_position($1::uuid[], id)
	`

	rows, err := bs.db.QueryContext(ctx, BookViewPreviewsQuery, pq.Array(bookIds))
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
		var currentObject domain.BookViewPreview
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.ImageUrl,
			&currentObject.Title,
			&currentObject.Author,
			&currentObject.CreatedDate,
			&currentObject.Genre,
			&currentObject.Rate,
			&currentObject.PagesCount,
		); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		currentObject.ViewedAt = viewedAtByBookId[currentObject.Id]
		previews = append(previews, currentObject)
	}

	if err := rows.Err(); err != nil {
		bs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	bs.l.Info("Successfully got book view previews")
	return previews, nil
}

func (bs *bookStorage) ReadingBookPreviews(ctx context.Context, lastReadingBooks []domain.LastReadingBook) ([]domain.ReadingBookPreview, error) {
	if len(lastReadingBooks) == 0 {
		return make([]domain.ReadingBookPreview, 0), nil
	}

	bookIds := make([]uuid.UUID, 0, len(lastReadingBooks))
	lastStartedAtByBookId := make(map[uuid.UUID]time.Time, len(lastReadingBooks))
	for _, lastReadingBook := range lastReadingBooks {
		bookIds = append(bookIds, lastReadingBook.BookId)
		lastStartedAtByBookId[lastReadingBook.BookId] = lastReadingBook.LastStartedAt
	}

	previews := make([]domain.ReadingBookPreview, 0, len(lastReadingBooks))
	const ReadingBookPreviewsQuery = `
		SELECT
			id,
			image_url,
			title,
			author,
			rate,
			genre,
			created_date,
			pages_count
		FROM books
		WHERE id = ANY($1)
		ORDER BY array_position($1::uuid[], id)
	`

	rows, err := bs.db.QueryContext(ctx, ReadingBookPreviewsQuery, pq.Array(bookIds))
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
		var currentObject domain.ReadingBookPreview
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.ImageUrl,
			&currentObject.Title,
			&currentObject.Author,
			&currentObject.Rate,
			&currentObject.Genre,
			&currentObject.CreatedDate,
			&currentObject.PagesCount,
		); err != nil {
			bs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		currentObject.LastStartedAt = lastStartedAtByBookId[currentObject.Id]
		previews = append(previews, currentObject)
	}

	if err := rows.Err(); err != nil {
		bs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	bs.l.Info("Successfully got reading book previews")
	return previews, nil
}
