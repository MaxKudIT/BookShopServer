package orders

import (
	"context"
	"errors"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (oserv *ordersService) Create(ctx context.Context, firebaseId string, order domain.Order) (uuid.UUID, error) {
	if order.TotalAmount <= 0 {
		err := errors.New("total amount must be greater than zero")
		oserv.l.Error("Invalid order total amount", "error", err)
		return uuid.Nil, err
	}

	userId, err := oserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		oserv.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}

	if order.Currency == "" {
		order.Currency = "RUB"
	}

	order.Id = uuid.New()
	order.UserId = userId
	order.Status = domain.OrderPaid
	order.PaidAt = time.Now()

	if err := oserv.os.Save(ctx, order); err != nil {
		oserv.l.Error("Error saving order", "error", err)
		return order.Id, err
	}

	oserv.l.Info("Successfully created order", "id", order.Id)
	return order.Id, nil
}
