package UserServices

import (
	"errors"
	"fmt"
	"shopperia/src/auth"
	"shopperia/src/common/models"
)

var condition string = ""

func (S *Service) LoginUser(form models.Login) (string, error) {
	if form.Email == condition || form.Password == condition {
		errStr := fmt.Sprintf("please provide all login fields")
		return "", errors.New(errStr)
	}

	ok, err := S.UseCase.CheckLogin(form.Email, form.Password)
	if err != nil || !ok {
		errStr := fmt.Sprintf("error while verifing data. err: %v", err)
		return "", errors.New(errStr)
	}

	userName, err := S.UseCase.GetUserNameByEmail(form.Email)
	if err != nil {
		errStr := fmt.Sprintf("Error while getting userName. ERR: %v", err)
		return "", errors.New(errStr)
	}

	token, err := auth.GenerateToken(form, userName, false)
	if err != nil {
		return "", errors.New("error while login")
	}

	return token, err
}
