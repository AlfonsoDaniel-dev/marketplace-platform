package UserServices

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	userModel "shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"time"
)

func (S *Service) Register(register UserDTO.RegisterDTO) error {
	if register.FirstName == "" || register.LastName == "" || register.Email == "" {
		return errors.New("please provide valid fields")
	}

	AccountCreatedAt := time.Now()
	UserId := uuid.New()

	user := userModel.User{
		Id:              UserId,
		FirstName:       register.FirstName,
		LastName:        register.LastName,
		UserName:        register.UserName,
		Biography:       register.Biography,
		ProfilePicture:  userModel.Image{},
		Age:             register.Age,
		Email:           register.Email,
		Password:        register.Password,
		UserAddress:     userModel.Address{},
		OrderedProducts: nil,
		Orders:          nil,
		CreatedAt:       AccountCreatedAt,
		UpdatedAt:       time.Time{},
	}

	fmt.Println(user.Biography)

	err := S.UseCase.CreateUserWithOutAddress(user)
	if err != nil {
		errString := fmt.Sprintf("rror while registering user: %v", err)
		return errors.New(errString)
	}
	/*
		err = S.UseCase.WelcomeEmail(user.FirstName, user.LastName, user.Email)
		if err != nil {
			fmt.Println(err)
			errStr := fmt.Sprintf("failed to send welcome email but register was success: %v", err)
			return errors.New(errStr)
		} */

	return nil
}
