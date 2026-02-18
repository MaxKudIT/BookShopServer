package cart_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (ciserv *cartItemService) IsInCart(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
		return false, err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)
	if err != nil {
		ciserv.l.Info("Failed to get cart id", "id", cartId)
		return false, err
	}

	IsInCart, err := ciserv.cis.IsInCart(ctx, cartId, bookId)
	if err != nil {
		ciserv.l.Error("result about book in the cart failed", "err", err)
		return false, err
	}

	ciserv.l.Info("get a result book in thr cart success")
	return IsInCart, nil
}

func (ciserv *cartItemService) AllCartItems(ctx context.Context, firebaseId string) ([]domain.CartItemPreview, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("cart saving failed", "err", err)
		return nil, err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)
	if err != nil {
		ciserv.l.Info("Failed to get cart id", "id", cartId)
		return nil, err
	}

	cartItems, err := ciserv.cis.AllCartItems(ctx, cartId)

	if err != nil {
		ciserv.l.Error("cart list failed", "err", err)
		return nil, err
	}

	ciserv.l.Info("cart list success")
	return cartItems, nil
}

func (ciserv *cartItemService) Create(ctx context.Context, firebaseId string, cartItem domain.CartItem) (uuid.UUID, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("cart saving failed", "err", err)
		return uuid.Nil, err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)
	if err != nil {
		ciserv.l.Info("Failed to create cart id", "id", cartId)
		return uuid.Nil, err
	}

	cartItem.CartId = cartId

	if err := ciserv.cis.Save(ctx, cartItem); err != nil {
		ciserv.l.Error("Error saving cart item", "error", err)
		return cartItem.BookId, err
	}
	ciserv.l.Info("Successfully created cart item")
	return cartItem.BookId, nil
}

func (ciserv *cartItemService) Delete(ctx context.Context, bookIds []uuid.UUID, firebaseId string) error {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("cart saving failed", "err", err)
		return err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)

	if err != nil {
		ciserv.l.Info("Failed to delete cart id", "id", cartId)
		return err
	}

	if err := ciserv.cis.Delete(ctx, bookIds, cartId); err != nil {
		ciserv.l.Error("Error deleting cart items", "error", err)
		return err
	}
	ciserv.l.Info("Successfully deleting cart items")
	return nil
}
