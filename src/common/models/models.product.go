package models

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	Id           uuid.UUID `json:"id"`
	CategoryName string    `json:"category_Name"`
	Description  string    `json:"description"`
}

type Product struct {
	ID              uuid.UUID `json:"id"`
	SellerID        uuid.UUID `json:"seller_id"`
	Name            string    `json:"name"`
	ProductCategory Category  `json:"product_category"`
	ProductPicture  Image     `json:"product_picture"`
	Description     string    `json:"description"`
	Stock           int64     `json:"stock"`
	Price           float32   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
