package UserServices

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"time"
)

func (S *Service) UploadAddress(email string, form UserDTO.UploadAddressForm) error {
	if email == "" || form.PostalCode == "" || form.Country == "" || form.City == "" || form.State == "" {
		return errors.New("please provide all fields required")
	}

	addressToUpdate := models.Address{
		ID:         uuid.UUID{},
		UserId:     uuid.UUID{},
		Street:     form.Street,
		City:       form.City,
		State:      form.State,
		PostalCode: form.Country,
		Country:    form.PostalCode,
		CreatedAt:  time.Time{},
	}

	err := S.UseCase.UploadAddressData(email, addressToUpdate)
	if err != nil {
		errStr := fmt.Sprintf("upload address data error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) UpdateUserName(form UserDTO.UpdateUserName) error {
	if form.NewUserName == "" || form.Common.Email == "" {
		return errors.New("please provide all fields required")
	}

	err := S.UseCase.ChangeUserName(form.NewUserName, form.Common.Email)
	if err != nil {
		errStr := fmt.Sprintf("change user name error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) UpdateUserFirstName(form UserDTO.UpdateFirstName) error {
	if form.FirstName == "" || form.Common.Email == "" {
		return errors.New("please provide all fields required")
	}

	err := S.UseCase.ChangeUserFirstName(form.FirstName, form.Common.Email)
	if err != nil {
		errStr := fmt.Sprintf("change user first name error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) UpdateUserLastName(form UserDTO.UpdateLastName) error {
	if form.LastName == "" || form.Common.Email == "" {
		return errors.New("please provide all fields required")
	}

	err := S.UseCase.ChangeUserLastName(form.LastName, form.Common.Email)
	if err != nil {
		errStr := fmt.Sprintf("change user last name error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) UpdateUserPassword(form UserDTO.UpdatePassword) error {
	if form.OldPassword == "" || form.NewPassword == "" || form.Common.Email == "" {
		errStr := fmt.Sprintf("please provide all fields required")
		return errors.New(errStr)
	}

	err := S.UseCase.ChangeUserPassword(form.Common.Email, form.OldPassword, form.NewPassword)
	if err != nil {
		errStr := fmt.Sprintf("change user password error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) UpdateUserEmail(form UserDTO.UpdateEmail) error {
	if form.NewEmail == "" || form.Common.Email == "" || form.Password == "" {
		errStr := fmt.Sprintf("please provide all fields required")
		return errors.New(errStr)
	}

	err := S.UseCase.ChangeUserEmail(form.NewEmail, form.Common.Email, form.Password)
	if err != nil {
		errStr := fmt.Sprintf("change user email error: %v", err)
		return errors.New(errStr)
	}
	
	return nil
}
