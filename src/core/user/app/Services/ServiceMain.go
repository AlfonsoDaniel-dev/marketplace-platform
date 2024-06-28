package UserServices

import (
	"shopperia/src/common/models"
	userApp "shopperia/src/core/user/app"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

type Service struct {
	userApp.UseCase
}

func NewService(useCase userApp.UseCase) *Service {
	return &Service{
		UseCase: useCase,
	}
}

type ServiceInterface interface {
	Hello() string
	Register(register UserDTO.RegisterDTO) error
	LoginUser(form models.Login) (string, error)
	CheckTwoStepsVerification(email string) (bool, error)
	SendLoginConfirmation(email string) (string, error)
	SendPasswordEmailConfirmation(email string) (string, error)
	CheckToken(email, token string) (bool, error)
	CleanToken(email string) error
	// upload Services

	UploadAddress(email string, form UserDTO.UploadAddressForm) error
	// update services

	UpdateUserName(form UserDTO.UpdateUserName) error
	UpdateUserFirstName(form UserDTO.UpdateFirstName) error
	UpdateUserLastName(form UserDTO.UpdateLastName) error
	UpdateUserPassword(form UserDTO.UpdatePassword) error
	UpdateUserEmail(form UserDTO.UpdateEmail) error
	UpdateUserTsvConfig(form UserDTO.UpdateTsvStatus) error
}

func (s *Service) Hello() string {
	return "Hello World"
}
