package user_domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/common/models"
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

func (u *UserDomain) CheckTwoStepsVerification(email string) (bool, error) {
	if email == nullString {
		return false, errors.New("please provide valid data")
	}

	ok, err := u.OutputInterface.PsqlCheckTwoStepsVerificationIsTrue(email)
	if err != nil || !ok {
		return false, err
	}

	return true, nil
}

func (u *UserDomain) GenerateAccessToken(email string) (string, error) {
	if email == nullString {
		return "", errors.New("please provide valid data")
	}

	codeBytes := make([]byte, 16)
	_, err := rand.Read(codeBytes)
	if err != nil {
		errStr := fmt.Sprintf("it was an error while generating access token: %v", err)
		return "", errors.New(errStr)
	}

	token := hex.EncodeToString(codeBytes)

	_, err = u.OutputInterface.PsqlInsertTsvCode(email, token)
	if err != nil {
		errStr := fmt.Sprintf("error while inserting access token: %v", err)
		return "", errors.New(errStr)
	}

	return token, nil
}

func (u *UserDomain) SendLoginConfirmationEmail(DestEmail string) (string, error) {
	if DestEmail == nullString {
		return "", errors.New("no email provide")
	}

	userName, err := u.OutputInterface.PsqlGetUserNameByEmail(DestEmail)
	if err != nil {
		return "", err
	}

	token, err := u.GenerateAccessToken(DestEmail)

	link := fmt.Sprintf("http://localhost:8080/api/user/login/confirm/%v/%v", DestEmail, token)

	text := "Hello there! Looks like you wanna log-in on your account, we just wanna make sure is You. Click bellow to confirm your trying to log-in."
	Linktext := "Click here to confirm"

	dataToSend := models.SendTSVLoginEmail{
		UserName:  userName,
		Text:      text,
		Link:      link,
		FinalText: Linktext,
	}

	if err := u.EmailInterface.SendLoginConfirmationEmail(dataToSend, DestEmail, userName); err != nil {
		fmt.Println(err)
		return "", err
	}

	return link, nil
}

func (u *UserDomain) CheckAccessToken(email, token string) (bool, error) {
	if email == nullString || token == nullString {
		return false, errors.New("please provide valid data")
	}

	token, err := u.OutputInterface.PsqlGetAccessToken(email)
	if err != nil {
		return false, err
	}

	if token == "" {
		return false, nil
	}

	return true, nil
}

func (u *UserDomain) CleanAccessToken(email string) error {
	if email == nullString {
		return errors.New("please provide valid data")
	}

	err := u.OutputInterface.PsqlCleanAccessToken(email)
	if err != nil {
		return err
	}

	return nil
}
