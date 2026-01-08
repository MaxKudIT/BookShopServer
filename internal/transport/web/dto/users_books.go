package dto

import "github.com/google/uuid"

type UsersBooksDTO struct {
	FirebaseId string
	BookId     uuid.UUID
}
