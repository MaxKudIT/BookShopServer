package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAdminAccessDenied = errors.New("admin access denied")
	ErrInvalidAdminBook  = errors.New("invalid admin book payload")
)

type AdminBookCreate struct {
	Title       string
	Author      string
	Genre       string
	Price       float64
	Discount    int
	ImageUrl    string
	Description string
	AboutBook   string
	Quote       string
	ReadingTime string
	Pages       []string
}

type AdminBookCreateResult struct {
	BookId     uuid.UUID
	PagesCount int
}
