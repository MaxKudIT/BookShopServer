package physical_books

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (pbs *physicalBooksStorage) All(ctx context.Context) ([]domain.PhysicalBook, error) {
	physicalBooks := make([]domain.PhysicalBook, 0)
	const AllPhysicalBooksQuery = `
		SELECT id, book_id, price, discount, format, stock_count
		FROM physical_books
		ORDER BY id
	`

	rows, err := pbs.db.QueryContext(ctx, AllPhysicalBooksQuery)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			pbs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			pbs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			pbs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.PhysicalBook
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.BookId,
			&currentObject.Price,
			&currentObject.Discount,
			&currentObject.Format,
			&currentObject.StockCount,
		); err != nil {
			pbs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		physicalBooks = append(physicalBooks, currentObject)
	}

	if err := rows.Err(); err != nil {
		pbs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	pbs.l.Info("Successfully got physical books")
	return physicalBooks, nil
}

func (pbs *physicalBooksStorage) ById(ctx context.Context, id uuid.UUID) (domain.PhysicalBook, error) {
	var physicalBook domain.PhysicalBook
	const PhysicalBookByIdQuery = `
		SELECT id, book_id, price, discount, format, stock_count
		FROM physical_books
		WHERE id = $1
	`

	if err := pbs.db.QueryRowContext(ctx, PhysicalBookByIdQuery, id).Scan(
		&physicalBook.Id,
		&physicalBook.BookId,
		&physicalBook.Price,
		&physicalBook.Discount,
		&physicalBook.Format,
		&physicalBook.StockCount,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			pbs.l.Error("Physical book not found", "error", err)
			return domain.PhysicalBook{}, err
		case errors.Is(err, context.Canceled):
			pbs.l.Warn("Query cancelled", "error", err)
			return domain.PhysicalBook{}, err
		case errors.Is(err, context.DeadlineExceeded):
			pbs.l.Warn("Query timed out", "error", err)
			return domain.PhysicalBook{}, err
		default:
			pbs.l.Error("Query failed", "error", err)
			return domain.PhysicalBook{}, err
		}
	}

	pbs.l.Info("Successfully got physical book", "id", id)
	return physicalBook, nil
}

func (pbs *physicalBooksStorage) IsPhysicalBookInStock(ctx context.Context, bookId uuid.UUID) (domain.PhysicalBookStockInfo, error) {
	physicalBookStockInfo := domain.PhysicalBookStockInfo{
		BookId:        bookId,
		PhysicalBooks: make([]domain.PhysicalBook, 0),
	}
	const PhysicalBookStockInfoQuery = `
		SELECT
			b.id,
			b.title,
			b.author,
			b.rate,
			pb.id,
			pb.book_id,
			pb.price,
			pb.discount,
			pb.format,
			pb.stock_count
		FROM books b
		LEFT JOIN physical_books pb ON pb.book_id = b.id
		WHERE b.id = $1
		ORDER BY pb.price
	`

	rows, err := pbs.db.QueryContext(ctx, PhysicalBookStockInfoQuery, bookId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			pbs.l.Warn("Query cancelled", "error", err)
			return domain.PhysicalBookStockInfo{}, err
		case errors.Is(err, context.DeadlineExceeded):
			pbs.l.Warn("Query timed out", "error", err)
			return domain.PhysicalBookStockInfo{}, err
		default:
			pbs.l.Error("Query failed", "error", err)
			return domain.PhysicalBookStockInfo{}, err
		}
	}
	defer rows.Close()

	foundBook := false
	for rows.Next() {
		foundBook = true

		var physicalBookId uuid.NullUUID
		var physicalBookBookId uuid.NullUUID
		var price sql.NullFloat64
		var discount sql.NullInt64
		var format sql.NullString
		var stockCount sql.NullInt64

		if err := rows.Scan(
			&physicalBookStockInfo.BookId,
			&physicalBookStockInfo.Title,
			&physicalBookStockInfo.Author,
			&physicalBookStockInfo.Rate,
			&physicalBookId,
			&physicalBookBookId,
			&price,
			&discount,
			&format,
			&stockCount,
		); err != nil {
			pbs.l.Error("Scan failed", "error", err)
			return domain.PhysicalBookStockInfo{}, err
		}

		if !physicalBookId.Valid {
			continue
		}

		physicalBook := domain.PhysicalBook{
			Id:         physicalBookId.UUID,
			BookId:     physicalBookBookId.UUID,
			Price:      price.Float64,
			Discount:   int(discount.Int64),
			Format:     format.String,
			StockCount: int(stockCount.Int64),
		}
		if physicalBook.StockCount > 0 {
			physicalBookStockInfo.IsInStock = true
		}

		physicalBookStockInfo.PhysicalBooks = append(physicalBookStockInfo.PhysicalBooks, physicalBook)
	}

	if err := rows.Err(); err != nil {
		pbs.l.Error("Rows failed", "error", err)
		return domain.PhysicalBookStockInfo{}, err
	}

	if !foundBook {
		pbs.l.Error("Book not found", "bookId", bookId)
		return domain.PhysicalBookStockInfo{}, sql.ErrNoRows
	}

	pbs.l.Info("Successfully got physical book stock info", "bookId", bookId)
	return physicalBookStockInfo, nil
}
