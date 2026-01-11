package book

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (b *bookService) AllBooks(ctx context.Context, firebaseId string) ([]domain.BookPreview, error) {

	userId, err := b.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		b.l.Error("user list failed", "err", err)
		return nil, err
	}
	books, err := b.bs.AllBooks(ctx, userId)
	if err != nil {
		b.l.Error("book list failed", "err", err)
		return nil, err
	}

	b.l.Info("book list success")
	return books, nil
}

func (b *bookService) AllMyBooks(ctx context.Context, firebaseId string) ([]domain.BookPreview, error) {

	userId, err := b.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		b.l.Error("user list failed", "err", err)
		return nil, err
	}

	books, err := b.bs.AllMyBooks(ctx, userId)

	if err != nil {
		b.l.Error("book list failed", "err", err)
		return nil, err
	}

	b.l.Info("book list success")
	return books, nil
}

func (b *bookService) BookById(ctx context.Context, firebaseId string, bookId uuid.UUID) (domain.Book, error) {
	userId, err := b.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		b.l.Error("user list failed", "err", err)
		return domain.Book{}, err
	}
	book, err := b.bs.BookById(ctx, userId, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b.l.Info("book not found", "error", err)
			return domain.Book{}, err
		}
		b.l.Info("Error getting book", "error", err)
		return domain.Book{}, err
	}

	b.l.Info("Successfully got book", "id", bookId)
	return book, nil
}

func (b *bookService) IsMyBook(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error) {
	userId, err := b.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		b.l.Error("user list failed", "err", err)
		return false, err
	}
	isMy, err := b.bs.IsMyBook(ctx, userId, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b.l.Info("book not found", "error", err)
			return false, err
		}
		b.l.Info("Error getting book", "error", err)
		return false, err
	}

	b.l.Info("Successfully got result about purchase", "id", bookId)
	return isMy, nil
}
