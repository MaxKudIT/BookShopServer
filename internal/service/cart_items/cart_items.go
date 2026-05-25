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

func (ciserv *cartItemService) Count(ctx context.Context, firebaseId string) (int, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
		return 0, err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)
	if err != nil {
		ciserv.l.Info("Failed to get cart id", "id", cartId)
		return 0, err
	}

	count, err := ciserv.cis.Count(ctx, cartId)
	if err != nil {
		ciserv.l.Error("result about count in the cart failed", "err", err)
		return 0, err
	}

	ciserv.l.Info("get a result count in thr cart success")
	return count, nil
}

func (ciserv *cartItemService) AreAllInCart(ctx context.Context, firebaseId string, bookIds []uuid.UUID) (bool, error) {

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

	IsInCart, err := ciserv.cis.AreAllInCart(ctx, cartId, bookIds)
	if err != nil {
		ciserv.l.Error("result about books in the cart failed", "err", err)
		return false, err
	}

	ciserv.l.Info("get a result books in thr cart success")
	return IsInCart, nil
}

func (ciserv *cartItemService) AllCartItems(ctx context.Context, firebaseId string) ([]domain.CartItemPreview, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
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
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
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

func (ciserv *cartItemService) CreateItems(ctx context.Context, firebaseId string, cartItems []domain.CartItem) (uuid.UUID, error) {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
		return uuid.Nil, err
	}

	cartId, err := ciserv.cs.CartByUserId(ctx, userId)
	if err != nil {
		ciserv.l.Info("Failed to create cart id", "id", cartId)
		return uuid.Nil, err
	}

	for i := range cartItems {
		cartItems[i].CartId = cartId
	}

	cartItemsId, err := ciserv.cis.AllCartItemsId(ctx, cartId)
	if err != nil {
		ciserv.l.Error("result about ids in the cart failed", "err", err)
		return uuid.Nil, err
	}

	var filteredItems []domain.CartItem = make([]domain.CartItem, 0)
	var isExistsMap = make(map[uuid.UUID]bool)

	for _, id := range cartItemsId {
		isExistsMap[id] = true
	}

	for _, item := range cartItems {
		if !isExistsMap[item.BookId] {
			filteredItems = append(filteredItems, item)
		}
	}

	if err := ciserv.cis.SaveFromFavs(ctx, filteredItems); err != nil {
		ciserv.l.Error("Error saving cart items", "error", err)
		return cartId, err
	}
	ciserv.l.Info("Successfully created cart items")
	return cartId, nil
}

func (ciserv *cartItemService) Delete(ctx context.Context, bookIds []uuid.UUID, firebaseId string) error {

	userId, err := ciserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ciserv.l.Error("failed getting id by firebaseId", "err", err)
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
