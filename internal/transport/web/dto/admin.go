package dto

import "github.com/google/uuid"

type AdminBookCreateDTO struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Genre       string   `json:"genre"`
	Price       float64  `json:"price"`
	Discount    int      `json:"discount"`
	ImageUrl    string   `json:"imageUrl"`
	Description string   `json:"description"`
	AboutBook   string   `json:"aboutBook"`
	Quote       string   `json:"quote"`
	ReadingTime string   `json:"readingTime"`
	Pages       []string `json:"pages"`
}

type AdminBookCreateResponseDTO struct {
	BookId     uuid.UUID `json:"bookId"`
	PagesCount int       `json:"pagesCount"`
}
