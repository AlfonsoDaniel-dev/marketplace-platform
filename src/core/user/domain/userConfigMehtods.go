package user_domain

import (
	"errors"
	"github.com/google/uuid"
	"shopperia/src/common/models"
	"shopperia/src/core/helpers"
	"time"
)

func (U *UserDomain) CreateUserWithOutAddress(user models.User) error {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" || user.UserName == "" {
		return errors.New("please provide all fields required")
	}

	encryptedPassword, err := helpers.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword

	err = U.OutputInterface.PsqlCreateUserWithOutAddress(user)
	if err != nil {
		return err
	}

	return nil
}

var null string = ""

func (u *UserDomain) UploadAddressData(email string, address models.Address) error {
	if email == null || address.City == null || address.State == null || address.Street == null || address.Country == null || address.PostalCode == null {
		return errors.New("please provide information for address fields")
	}

	address.ID = uuid.New()
	address.CreatedAt = time.Now()

	userID, err := u.PsqlGetUserIdByEmail(email)
	if err != nil {
		return err
	}

	address.UserId = userID

	return nil
}
