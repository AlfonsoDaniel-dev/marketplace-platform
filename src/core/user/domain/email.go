package user_domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"shopperia/src/common/models"
)

var nullstring string = ""

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

	host := os.Getenv("APP_HOST")

	link := fmt.Sprintf("http://%v:8080/api/user/login/confirm/%v/%v", host, DestEmail, token)

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

func (u *UserDomain) SendPasswordConfirmationEmail(DestEmail string) (string, error) {
	if DestEmail == nullString {
		return "", errors.New("please provide valid data")
	}

	token, err := u.GenerateAccessToken(DestEmail)
	if err != nil {
		return "", err
	}

	userName, err := u.OutputInterface.PsqlGetUserNameByEmail(DestEmail)
	if err != nil {
		return "", err
	}

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")

	link := fmt.Sprintf("http://%v:%v/api/v1/user/config/update/password/confirm/%v/%v", HOST, PORT, DestEmail, token)

	content := models.PasswordChangeEmail{
		UserName:         userName,
		ConfirmationLink: link,
	}

	err = u.EmailInterface.SendPasswordChangeConfirmationEmail(content, DestEmail, userName)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (u *UserDomain) SendTsvChangesConfirmationEmail(DestEmail string) (string, error) {
	if DestEmail == nullString {
		return "", errors.New("please provide valid data")
	}

	userName, err := u.OutputInterface.PsqlGetUserNameByEmail(DestEmail)
	if err != nil {
		return "", err
	}
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	token, err := u.GenerateAccessToken(DestEmail)
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf("http://%v:%v/api/user/config/tsv/confirm/%v/%v", host, port, DestEmail, token)

	content := models.TsvChangeEmail{
		UserName: userName,
		Link:     link,
	}

	err = u.EmailInterface.SendTsvChangeConfirmation(content, DestEmail, userName)
	if err != nil {
		return "", err
	}

	return link, nil
}

const welcomeText = "Hello!, yo just create an account on our platform we wish you to have a great experience with us."

func (u *UserDomain) WelcomeEmail(firstName, lastName, DestEmail string) error {
	if firstName == "" || lastName == "" || DestEmail == "" {
		return errors.New("all fields are required")
	}

	subject := "Registro en shopperia"

	usermame := firstName + " " + lastName
	title := "Welcome To Shopperia"

	content := models.WelcomeEmail{
		UserName:    usermame,
		Title:       title,
		WelcomeText: welcomeText,
	}

	mailInfo := models.EmailDto{
		Subject:          subject,
		DestinationName:  usermame,
		DestinationEmail: DestEmail,
	}

	err := u.EmailInterface.SendWelcomeEmail(content, mailInfo)
	if err != nil {
		return err
	}

	return nil
}
