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

	// userUpdates

	ChangeUserName(NewUserName, email string) error
	ChangeUserFirstName(NewUserFirstName, email string) error
	ChangeUserLastName(NewUserLastName, email string) error
	ChangeUserEmail(NewUserEmail, email, password string) error
	ChangeUserPassword(email, oldpassword, newPassword string) error
	ChangeUserTsvStatus(email string, value bool) error

	// userGetters

	GetUserIdByEmail(email string) (uuid.UUID, error)
	CheckLogin(email, password string) (bool, error)
	CheckTwoStepsVerification(email string) (bool, error)
	SendLoginConfirmationEmail(DestEmail string) (string, error)
	SendPasswordConfirmationEmail(DestEmail string) (string, error)
	SendTsvChangesConfirmationEmail(DestEmail string) (string, error)
	CheckAccessToken(email, token string) (bool, error)
	CleanAccessToken(email string) error

	WelcomeEmail(firstName, lastName, DestEmail string) error
}
