package user

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (u *userService) Create(ctx context.Context, userp domain.User) (domain.User, error) { //DOMAIN

	if err := u.us.Save(ctx, userp); err != nil {
		u.l.Error("Error saving user", "error", err)
		return domain.User{}, err
	}
	u.l.Info("Successfully saving user")
	return userp, nil
}

func (u *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.us.Delete(ctx, id); err != nil {
		u.l.Error("Error deleting user", "error", err)
		return err
	}
	u.l.Info("Successfully deleting user")
	return nil
}
