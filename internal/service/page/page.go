package page

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (p *pageService) AllPagesOfBook(ctx context.Context, bookId uuid.UUID) (int, error) {
	pages, err := p.ps.AllPagesOfBook(ctx, bookId)

	if err != nil {
		p.l.Error("pages count failed", "err", err)
		return 0, err
	}

	p.l.Info("pages count success")
	return pages, nil
}

func (p *pageService) PageByNumber(ctx context.Context, pageNumber int, bookId uuid.UUID) (domain.Page, error) {
	page, err := p.ps.PageByNumber(ctx, pageNumber, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.l.Info("page not found", "error", err)
			return domain.Page{}, err
		}
		p.l.Info("Error getting page", "error", err)
		return domain.Page{}, err
	}

	p.l.Info("Successfully got page", "id", page.Id)
	return page, nil
}
