package user_domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/core/helpers"
)

var nullString string = ""

func (u *UserDomain) GetUserIdByEmail(email string) (uuid.UUID, error) {
	if email == "" {
		return uuid.Nil, errors.New("pleaseProvide a valid email")
	}

	id, err := u.PsqlGetUserIdByEmail(email)
	if err != nil {
		errStr := fmt.Sprintf("error Getting string: %v", err)
		return uuid.Nil, errors.New(errStr)
	}

	return id, nil
}

func (u *UserDomain) GetUserNameByEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("search email is empty, need email for searching")
	}

	userName, err := u.OutputInterface.PsqlGetUserNameByEmail(email)
	if err != nil {
		return "", err
	}

	return userName, nil
}

func (u *UserDomain) CheckLogin(email, password string) (bool, error) {
	if email == nullString || password == nullString {
		return false, errors.New("please provide valid login form data")
	}

	dbEmail, err := u.OutputInterface.PsqlVerifyEmailExists(email)
	if err != nil {
		errStr := fmt.Sprintf("error while verifying email %v, may be email isn't valid", err)
		return false, errors.New(errStr)
	}

	hashedPassword, err := u.OutputInterface.PsqlGetHashPassword(dbEmail)
	if err != nil {
		errStr := fmt.Sprintf("error while verifying password %v", err)
		return false, errors.New(errStr)
	}

	ok := helpers.ComparePasswords(hashedPassword, password)
	if !ok {
		return false, errors.New("error")
	}

	return true, nil
}
