package user_domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	err = u.PsqlInsertAddressData(address)

	return nil
}

func (u *UserDomain) ChangeUserName(NewUserName, email string) error {
	if NewUserName == null || email == null {
		return errors.New("please provide information for user fields")
	}

	err := u.OutputInterface.PsqlChangeUserName(NewUserName, email)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDomain) ChangeUserFirstName(NewFirstName, email string) error {
	if NewFirstName == null || email == null {
		return errors.New("please provide information for user fields")
	}

	err := u.OutputInterface.PsqlChangeUserFirstName(NewFirstName, email)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDomain) ChangeUserLastName(NewLastName, email string) error {
	if NewLastName == null || email == null {
		return errors.New("please provide information for user fields")
	}

	err := u.OutputInterface.PsqlChangeUserLastName(NewLastName, email)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDomain) validatePassword(password, email string) (bool, error) {
	if password == null || email == "" {
		return false, errors.New("please provide information for password fields")
	}

	hash, err := u.OutputInterface.PsqlGetHashPassword(email)
	if err != nil {
		return false, err
	}

	ok := bcrypt.CompareHashAndPassword(hash, hash)
	if ok != nil {
		errStr := fmt.Sprintf("Error while checking password. ERR: %v", err)
		return false, errors.New(errStr)
	}

	return true, nil
}

func (u *UserDomain) ChangeUserEmail(NewEmail, email, password string) error {
	if NewEmail == null || email == null || password == null {
		return errors.New("please provide information for user fields")
	}

	ok, err := u.validatePassword(password, email)
	if err != nil || !ok {
		errStr := fmt.Sprintf("Error while validating password. may be password is wrong ERR: %v", err)
		return errors.New(errStr)
	}

	err = u.OutputInterface.PsqlChangeUserEmail(NewEmail, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDomain) ChangeUserPassword(email, oldpassword, newPassword string) error {
	if email == null || oldpassword == null || newPassword == null {
		return errors.New("please provide information for password fields")
	}

	ok, err := u.validatePassword(oldpassword, email)
	if err != nil || !ok {
		errStr := fmt.Sprintf("Error while validating password. may be password is wrong ERR: %v", err)
		return errors.New(errStr)
	}

	err = u.OutputInterface.PsqlChangeUserPassword(email, oldpassword, newPassword)
	if err != nil {
		return err
	}

	return nil
}
