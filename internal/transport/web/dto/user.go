package dto

import (
	"github.com/bookshop/internal/domain"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type UserDTO struct {
	Login    string
	Password string
	Email    string
}

func Validate(u UserDTO) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Login,
			validation.Required,
			validation.Length(3, 20),
		),
		validation.Field(&u.Password,
			validation.Required,
			validation.Length(8, 100),
		),
		validation.Field(&u.Email,
			validation.Required,
			is.Email,
		),
	)
}

func UserToDomain(id uuid.UUID, passwordHash string, userDTO UserDTO) domain.User {

	return domain.User{
		Id:           id,
		Login:        userDTO.Login,
		PasswordHash: passwordHash,
		Email:        userDTO.Email,
	}

}
