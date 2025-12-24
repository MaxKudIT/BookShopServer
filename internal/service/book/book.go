package book

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (b *bookService) AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {
	books, err := b.bs.AllBooks(ctx, userId)

	if err != nil {
		b.l.Error("book list failed", "err", err)
		return nil, err
	}

	b.l.Info("book list success")
	return books, nil
}

func (b *bookService) AllMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error) {
	books, err := b.bs.AllMyBooks(ctx, userId)

	if err != nil {
		b.l.Error("book list failed", "err", err)
		return nil, err
	}

	b.l.Info("book list success")
	return books, nil
}

func (b *bookService) BookById(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	book, err := b.bs.BookById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b.l.Info("book not found", "error", err)
			return domain.Book{}, err
		}
		b.l.Info("Error getting book", "error", err)
		return domain.Book{}, err
	}

	b.l.Info("Successfully got book", "id", id)
	return book, nil
}
