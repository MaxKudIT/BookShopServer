package order_items

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (oiserv *orderItemsService) Create(ctx context.Context, firebaseId string, orderItem domain.OrderItem) error {
	if orderItem.Quantity <= 0 {
		err := errors.New("quantity must be greater than zero")
		oiserv.l.Error("Invalid order item quantity", "error", err)
		return err
	}

	userId, err := oiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		oiserv.l.Error("Error getting userId by firebaseId", "error", err)
		return err
	}

	if err := oiserv.ois.Save(ctx, userId, orderItem); err != nil {
		oiserv.l.Error("Error saving order item", "error", err)
		return err
	}

	oiserv.l.Info("Successfully created order item", "orderId", orderItem.OrderId)
	return nil
}

func (oiserv *orderItemsService) All(ctx context.Context, firebaseId string, orderId uuid.UUID) ([]domain.OrderItem, error) {
	userId, err := oiserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		oiserv.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	orderItems, err := oiserv.ois.AllByOrderId(ctx, userId, orderId)
	if err != nil {
		oiserv.l.Error("Error getting order items", "error", err)
		return nil, err
	}

	oiserv.l.Info("Successfully got order items", "orderId", orderId)
	return orderItems, nil
}
