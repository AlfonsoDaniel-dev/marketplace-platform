package models

import (
	"github.com/google/uuid"
	"time"
)

type UserShopOrder struct {
	ID       uuid.UUID    `json:"id"`
	SellerID uuid.UUID    `json:"seller_id"`
	Items    []OrderItems `json:"items"`
	Total    float64      `json:"total"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItems struct {
	ID       uuid.UUID `json:"id"`
	OrderID  uuid.UUID `json:"order_id"`
	Product  Product   `json:"product_id"`
	Quantity int64     `json:"quantity"`
	Total    float64   `json:"total"`
}

type PurchaseProduct struct {
	OrderId  uuid.UUID `json:"order_id"`
	BuyerID  uuid.UUID `json:"buyer_id"`
	SellerID uuid.UUID `json:"seller_id"`
	Product  Product   `json:"product"`
	Quantity float64   `json:"quantity"`
	Total    float64   `json:"total"`
}
