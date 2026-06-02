package dto

import "github.com/bookshop/internal/domain"

type OrderDTO struct {
	TotalAmount     float64 `json:"totalAmount"`
	Currency        string  `json:"currency"`
	DeliveryAddress string  `json:"deliveryAddress"`
}

func OrderToDomain(orderDTO OrderDTO) domain.Order {
	return domain.Order{
		TotalAmount:     orderDTO.TotalAmount,
		Currency:        orderDTO.Currency,
		DeliveryAddress: orderDTO.DeliveryAddress,
	}
}
