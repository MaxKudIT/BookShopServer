package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type UserDTO struct {
	FirebaseId string
}

func UserToDomain(id uuid.UUID, userDTO UserDTO) domain.User {

	return domain.User{
		Id:         id,
		FirebaseId: userDTO.FirebaseId,
	}

}
