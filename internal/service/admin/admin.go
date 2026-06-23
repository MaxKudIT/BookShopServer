package admin

import (
	"context"
	"strings"

	"github.com/bookshop/internal/domain"
)

func (aserv *adminService) CreateBook(ctx context.Context, firebaseId string, book domain.AdminBookCreate) (domain.AdminBookCreateResult, error) {
	isAdmin, err := aserv.as.IsAdmin(ctx, firebaseId)
	if err != nil {
		aserv.l.Error("failed to check admin access", "error", err)
		return domain.AdminBookCreateResult{}, err
	}
	if !isAdmin {
		aserv.l.Warn("admin access denied", "firebaseId", firebaseId)
		return domain.AdminBookCreateResult{}, domain.ErrAdminAccessDenied
	}

	normalizedBook, err := normalizeAdminBook(book)
	if err != nil {
		aserv.l.Warn("invalid admin book payload", "error", err)
		return domain.AdminBookCreateResult{}, err
	}

	return aserv.as.CreateBookWithPages(ctx, normalizedBook)
}

func normalizeAdminBook(book domain.AdminBookCreate) (domain.AdminBookCreate, error) {
	book.Title = strings.TrimSpace(book.Title)
	book.Author = strings.TrimSpace(book.Author)
	book.Genre = strings.TrimSpace(book.Genre)
	book.ImageUrl = strings.TrimSpace(book.ImageUrl)
	book.Description = strings.TrimSpace(book.Description)
	book.AboutBook = strings.TrimSpace(book.AboutBook)
	book.Quote = strings.TrimSpace(book.Quote)
	book.ReadingTime = strings.TrimSpace(book.ReadingTime)

	if book.Title == "" || book.Author == "" || book.Genre == "" {
		return domain.AdminBookCreate{}, domain.ErrInvalidAdminBook
	}
	if book.Price < 0 || book.Discount < 0 || book.Discount > 100 {
		return domain.AdminBookCreate{}, domain.ErrInvalidAdminBook
	}
	if len(book.Pages) == 0 {
		return domain.AdminBookCreate{}, domain.ErrInvalidAdminBook
	}

	pages := make([]string, 0, len(book.Pages))
	for _, page := range book.Pages {
		page = strings.TrimSpace(page)
		if page == "" {
			continue
		}
		pages = append(pages, page)
	}
	if len(pages) == 0 {
		return domain.AdminBookCreate{}, domain.ErrInvalidAdminBook
	}

	book.Pages = pages
	return book, nil
}
