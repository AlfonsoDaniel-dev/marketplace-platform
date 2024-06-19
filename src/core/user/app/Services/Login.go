package UserServices

import (
	"errors"
	"fmt"
	"shopperia/src/auth"
	"shopperia/src/common/models"
)

var condition string = ""

func (S *Service) SendConfirmationEmain(email string) (string, error) {
	if email == condition {
		return "", nil
	}

	link, err := S.UseCase.SendLoginConfirmationEmail(email)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (S *Service) CheckLoginData(form models.Login) (bool, error) {
	if form.Email == "" {
		errStr := fmt.Sprintf("email cannot be empty")
		return false, errors.New(errStr)
	}

	ok, err := S.UseCase.CheckLogin(form.Email, form.Password)
	if err != nil || !ok {
		errStr := fmt.Sprintf("it was an error checking login information, may be email or password ar null")
		return false, errors.New(errStr)
	}

	return true, nil
}

func (S *Service) LoginUser(form models.Login) (string, error) {
	if form.Email == condition || form.Password == condition {
		errStr := fmt.Sprintf("please provide all login fields")
		return "", errors.New(errStr)
	}

	userName, err := S.UseCase.GetUserNameByEmail(form.Email)
	if err != nil {
		errStr := fmt.Sprintf("Error while getting userName. ERR: %v", err)
		return "", errors.New(errStr)
	}

	token, err := auth.GenerateToken(form, userName, false)
	if err != nil {
		errStr := fmt.Sprintf("Error while generating token. ERR: %v", err)
		return "", errors.New(errStr)
	}

	return token, nil
}

func (S *Service) CheckTwoStepsVerification(email string) (bool, error) {
	if email == "" {
		errStr := fmt.Sprintf("No fields provide")
		return false, errors.New(errStr)
	}

	ok, err := S.UseCase.CheckTwoStepsVerification(email)
	if !ok || err != nil {
		errStr := fmt.Sprintf("User Does not have Two Steps Verification or it was an error: %v", err)
		return false, errors.New(errStr)
	}

	return true, nil
}

func (S *Service) SendLoginConfirmation(email string) (string, error) {
	if email == "" {
		errStr := fmt.Sprintf("please provide all login fields")
		return "", errors.New(errStr)
	}

	ConfirmationLink, err := S.UseCase.SendLoginConfirmationEmail(email)
	if err != nil {
		errStr := fmt.Sprintf("it was an error while sendingConfirmation email: %v", err)
		return "", errors.New(errStr)
	}

	return ConfirmationLink, nil
}
