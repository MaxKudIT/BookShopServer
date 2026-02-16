package cart

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (cserv *cartService) Create(ctx context.Context, cart domain.Cart, firebaseId string) (uuid.UUID, error) {

	userId, err := cserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		cserv.l.Error("cart saving failed", "err", err)
		return uuid.Nil, err
	}
	cart.UserId = userId

	if err := cserv.cs.Save(ctx, cart); err != nil {
		cserv.l.Error("Error saving cart", "error", err)
		return cart.Id, err
	}
	cserv.l.Info("Successfully created cart")
	return cart.Id, nil
}
