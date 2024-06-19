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
	UploadAddress(email string, form UserDTO.UploadAddressForm) error
}

func (s *Service) Hello() string {
	return "Hello World"
}
