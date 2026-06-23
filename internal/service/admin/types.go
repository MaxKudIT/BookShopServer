package admin

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
)

type adminStorage interface {
	IsAdmin(ctx context.Context, firebaseId string) (bool, error)
	CreateBookWithPages(ctx context.Context, book domain.AdminBookCreate) (domain.AdminBookCreateResult, error)
}

type adminService struct {
	as adminStorage
	l  *slog.Logger
}

func New(as adminStorage, l *slog.Logger) *adminService {
	return &adminService{as: as, l: l}
}
