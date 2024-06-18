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
