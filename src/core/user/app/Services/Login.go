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
	if err != nil {
		errStr := fmt.Sprintf("it was an error checking login information, may be email or password ar null")
		return false, errors.New(errStr)
	} else if !ok {
		return false, nil
	}

	return true, nil
}

func (S *Service) LoginUser(form models.Login) models.LoginStatus {
	if form.Email == condition || form.Password == condition {
		errStr := fmt.Sprintf("please provide all login fields")
		return models.LoginStatus{Error: errors.New(errStr)}
	}

	loginProcces := models.LoginStatus{}

	ok, err := S.UseCase.CheckLogin(form.Email, form.Password)
	if err != nil {
		errStr := fmt.Sprintf("it was an error checking login information, may be email or password are bad, ERR %v", err)
		return models.LoginStatus{Error: errors.New(errStr)}
	} else if !ok {
		errStr := fmt.Sprintf("please provide all login fields")
		return models.LoginStatus{Error: errors.New(errStr)}
	}

	hasTSV, err := S.UseCase.CheckTwoStepsVerification(form.Email)
	if err != nil {
		return models.LoginStatus{Error: err}
	} else if !hasTSV {

		loginProcces.HasTwoStepsVerification = false
		userName, err := S.UseCase.GetUserNameByEmail(form.Email)
		if err != nil {
			errStr := fmt.Sprintf("Error while getting userName. ERR: %v", err)
			return models.LoginStatus{Error: errors.New(errStr)}
		}

		token, err := auth.GenerateToken(form, userName, false)
		if err != nil {
			errStr := fmt.Sprintf("Error while generating token. ERR: %v", err)
			return models.LoginStatus{Error: errors.New(errStr)}
		}

		loginProcces.Token = token

		return loginProcces
	}

	confimationLink, err := S.UseCase.SendLoginConfirmationEmail(form.Email)
	if err != nil {
		errStr := fmt.Sprintf("Error while sending login confirmation email. ERR: %v", err)
		return models.LoginStatus{Error: errors.New(errStr)}
	}

	loginProcces.HasTwoStepsVerification = true
	loginProcces.TsvConfirmationLink = confimationLink
	return loginProcces

	return loginProcces
}

func (S *Service) HandleTSVConfirmation(login models.Login, Accestoken string) (string, error) {
	tokenValid, err := S.UseCase.CheckAccessToken(login.Email, Accestoken)
	if err != nil {
		return "", err
	}

	if !tokenValid {
		return "", errors.New("token is invalid")
	}

	userName, err := S.UseCase.GetUserNameByEmail(login.Email)

	tokenjwt, err := auth.GenerateToken(login, userName, false)
	if err != nil {
		return "", err
	}

	err = S.UseCase.CleanAccessToken(login.Email)
	if err != nil {
		return "", err
	}

	return tokenjwt, nil
}

func (S *Service) CheckTwoStepsVerification(email string) (bool, error) {
	if email == "" {
		errStr := fmt.Sprintf("No fields provide")
		return false, errors.New(errStr)
	}

	ok, err := S.UseCase.CheckTwoStepsVerification(email)
	if err != nil {
		errStr := fmt.Sprintf("User Does not have Two Steps Verification or it was an error: %v", err)
		return false, errors.New(errStr)
	} else if !ok {
		return false, nil
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
		fmt.Println(err)
		errStr := fmt.Sprintf("it was an error while sendingConfirmation email: %v", err)
		return "", errors.New(errStr)
	}

	return ConfirmationLink, nil
}

func (S *Service) CheckToken(email, token string) (bool, error) {
	if email == "" || token == "" {
		return false, errors.New("email or token cannot be empty")
	}

	ok, err := S.UseCase.CheckAccessToken(email, token)
	if err != nil || !ok {
		errStr := fmt.Sprintf("It was an error verifing token, may be token is bad. ERR: %v", err)
		return false, errors.New(errStr)
	}

	return true, nil
}

func (S *Service) CleanToken(email string) error {
	if email == "" {
		return errors.New("no email provide")
	}

	err := S.UseCase.CleanAccessToken(email)
	fmt.Println(err)
	if err != nil {
		errStr := fmt.Sprintf("Error while cleaning token. ERR: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (S *Service) SendPasswordEmailConfirmation(email string) (string, error) {
	if email == "" {
		return "", errors.New("no email provide")
	}

	link, err := S.UseCase.SendPasswordConfirmationEmail(email)
	if err != nil {
		errStr := fmt.Sprintf("Error while sending password email confirmation. ERR: %v", err)
		return "", errors.New(errStr)
	}

	return link, nil
}

func (S *Service) SendTsvChangeEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("Please provide a valid email")
	}

	link, err := S.UseCase.SendTsvChangesConfirmationEmail(email)
	if err != nil {
		errStr := fmt.Sprintf("Error while sending tsv change email. ERR: %v", err)
		return "", errors.New(errStr)
	}

	return link, nil
}
