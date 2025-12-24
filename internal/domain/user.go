package domain

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Login        string
	Email        string
	PasswordHash string
}
