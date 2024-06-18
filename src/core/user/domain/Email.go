package user_domain

import (
	"errors"
	"shopperia/src/common/models"
)

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
