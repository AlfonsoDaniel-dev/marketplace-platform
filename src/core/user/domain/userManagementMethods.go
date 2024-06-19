package user_domain

import (
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

func (u *UserDomain) SendLoginConfirmationEmail(DestEmail string) (string, error) {
	if DestEmail == nullString {
		return "", errors.New("no email provide")
	}

	UserName, err := u.OutputInterface.PsqlGetUserNameByEmail(DestEmail)
	if err != nil {
		return "", err
	}

	type ConfirmationDataToSend struct {
		UserName  string
		Text      string
		Link      string
		FinalText string
	}

	link := fmt.Sprintf("/api/login/confirm/%v/%v", DestEmail, UserName)

	text := "Hello there! Looks like you wanna log-in on your account, we just wanna make sure is You. Click bellow to confirm your trying to log-in."
	Linktext := "Click here to confirm"

	dataToSend := ConfirmationDataToSend{
		UserName:  UserName,
		Text:      text,
		Link:      link,
		FinalText: Linktext,
	}

	EmailInfo := models.SendEmailForm{
		Subject:          "Login Confirmation",
		DestinationEmail: DestEmail,
		DestinationName:  UserName,
		TemplateData:     dataToSend,
	}

	if err := u.EmailInterface.SendLoginConfirmationEmail(EmailInfo); err != nil {
		return "", err
	}

	return link, nil
}
