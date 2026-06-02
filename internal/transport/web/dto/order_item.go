package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type OrderItemDTO struct {
	PhysicalProductId uuid.UUID `json:"physicalProductId"`
	Quantity          int       `json:"quantity"`
}

func OrderItemToDomain(orderId uuid.UUID, orderItemDTO OrderItemDTO) domain.OrderItem {
	return domain.OrderItem{
		OrderId:           orderId,
		PhysicalProductId: orderItemDTO.PhysicalProductId,
		Quantity:          orderItemDTO.Quantity,
	}
}
