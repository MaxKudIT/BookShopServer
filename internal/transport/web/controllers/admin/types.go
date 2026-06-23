package admin

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/gin-gonic/gin"
)

type adminService interface {
	CreateBook(ctx context.Context, firebaseId string, book domain.AdminBookCreate) (domain.AdminBookCreateResult, error)
}

type adminHandler struct {
	as adminService
	l  *slog.Logger
}

func New(as adminService, l *slog.Logger) *adminHandler {
	return &adminHandler{as: as, l: l}
}

type adminHandlers interface {
	CreateBook(ctx context.Context, c *gin.Context)
}
