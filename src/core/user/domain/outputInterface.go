package user_domain

import (
	"github.com/google/uuid"
	models "shopperia/src/common/models"
)

type OutputInterface interface {
	PsqlCreateUserWithOutAddress(user models.User) error
	PsqlGetUserIdByEmail(email string) (uuid.UUID, error)
	PsqlGetUserNameByEmail(email string) (string, error)
	PsqlInsertAddressData(address models.Address) error

	// login methods
	PsqlGetHashPassword(email string) ([]byte, error)
	PsqlVerifyEmailExists(email string) (string, error)
}

type EmailInterface interface {
	SendWelcomeEmail(emailContent models.WelcomeEmail, email models.EmailDto) error
	SendLoginConfirmationEmail(email models.EmailDto) error
}
