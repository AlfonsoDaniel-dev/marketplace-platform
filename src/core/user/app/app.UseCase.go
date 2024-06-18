package userApp

import (
	"github.com/google/uuid"
	userModel "shopperia/src/common/models"
)

// UseCase Uses methods on domain
type UseCase interface {
	// user config

	CreateUserWithOutAddress(user userModel.User) error
	UploadAddressData(email string, address userModel.Address) error
	GetUserNameByEmail(email string) (string, error)

	// userGetters

	GetUserIdByEmail(email string) (uuid.UUID, error)
	CheckLogin(email, password string) (bool, error)

	WelcomeEmail(firstName, lastName, DestEmail string) error
}
