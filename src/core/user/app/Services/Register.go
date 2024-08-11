package UserServices

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	userModel "shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"strings"
	"time"
)

func (S *Service) Register(register UserDTO.RegisterDTO) error {
	if register.FirstName == "" || register.LastName == "" || register.Email == "" {
		errStr := fmt.Sprintf("Email from form: %s, FirstName from form: %s, LastName from form: %s", register.Email, register.FirstName, register.LastName)
		return errors.New(errStr)
	}

	register.UserName = strings.ReplaceAll(register.UserName, " ", "_")
	register.Email = strings.ReplaceAll(register.Email, " ", "_")

	AccountCreatedAt := time.Now().Unix()
	UserId := uuid.New()

	user := userModel.User{
		Id:                  UserId,
		FirstName:           register.FirstName,
		LastName:            register.LastName,
		TwoStepsVerfication: register.TwoStepsVerification,
		UserName:            register.UserName,
		Biography:           register.Biography,
		ProfilePicture:      userModel.Image{},
		Age:                 register.Age,
		Email:               register.Email,
		Password:            register.Password,
		UserAddress:         userModel.Address{},
		OrderedProducts:     nil,
		Orders:              nil,
		CreatedAt:           AccountCreatedAt,
		UpdatedAt:           0,
	}

	fmt.Println(user.Biography)

	err := S.UseCase.CreateUserWithOutAddress(user)
	if err != nil {
		fmt.Println(err)
		errString := fmt.Sprintf("error while registering user: %v", err)
		return errors.New(errString)
	}

	err = S.UseCase.WelcomeEmail(user.FirstName, user.LastName, user.Email)
	if err != nil {
		fmt.Println(err)
		errStr := fmt.Sprintf("failed to send welcome email but register was success: %v", err)
		return errors.New(errStr)
	}

	return nil
}
