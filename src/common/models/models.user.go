package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id              uuid.UUID         `json:"id"`
	FirstName       string            `json:"first_Name"`
	LastName        string            `json:"last_Name"`
	UserName        string            `json:"user_name"`
	ProfilePicture  Image             `json:"profile_picture"`
	Collections     []Collection      `json:"collections"`
	Biography       string            `json:"biography"`
	Age             int               `json:"birth_date"`
	Email           string            `json:"email"`
	Password        string            `json:"password"`
	UserAddress     Address           `json:"user_address"`
	OrderedProducts []PurchaseProduct `json:"orders_created"`
	Orders          []UserShopOrder   `json:"orders"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
