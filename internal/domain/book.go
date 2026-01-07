package domain

import (
	"github.com/google/uuid"
	"time"
)

type Book struct {
	Id          uuid.UUID
	Title       string
	PagesCount  int
	Description string
	AboutBook   string
	Quote       string
	CreatedDate time.Time
	ReadingTime string
	Price       float64
	Discount    int //в процентах
	Author      string
	Genre       string
	ImageUrl    string
	Rate        float64
	IsMine      bool
}

type Genre string

const (
	Adv     Genre = "Приключения"
	Drama   Genre = "Драма"
	Horror  Genre = "Ужасы"
	History Genre = "Исторические"
	Fant    Genre = "Фантастика"
)

type BookPreview struct {
	Id       uuid.UUID
	Title    string
	Genre    Genre
	Price    float64
	Discount int
	ImageUrl string
	IsMine   bool
}
