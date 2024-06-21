package Userstorage

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	user_model "shopperia/src/common/models"
	"shopperia/src/core/helpers"
)

type psqlUser struct {
	DB *sql.DB
}

func NewPsqlUser(db *sql.DB) *psqlUser {
	return &psqlUser{
		DB: db,
	}
}

func (p *psqlUser) PsqlCreateUserWithOutAddress(user user_model.User) error {
	stmt, err := p.DB.Prepare(sqlCreateUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	nullTime := helpers.IntToNull(user.UpdatedAt)

	nullBiography := helpers.StringToNull(user.Biography)

	_, err = stmt.Exec(user.Id, user.FirstName, user.LastName, user.UserName, nullBiography, user.Age, user.Email, user.Password, user.TwoStepsVerfication, user.CreatedAt, nullTime)
	if err != nil {
		return err
	}

	log.Println("User created successfully")
	return nil
}

func (p *psqlUser) PsqlGetUserIdByEmail(email string) (uuid.UUID, error) {
	stmt, err := p.DB.Prepare(sqlGetUserIdByEmail)
	if err != nil {
		return uuid.Nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)
	var id uuid.UUID
	err = row.Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (p *psqlUser) PsqlInsertAddressData(address user_model.Address) error {
	stmt, err := p.DB.Prepare(sqlInsertUserAddress)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(address.ID, address.UserId, address.Street, address.City, address.State, address.PostalCode, address.Country, address.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) PsqlGetUserNameByEmail(email string) (string, error) {
	stmt, err := p.DB.Prepare(sqlGetUserNameByEmail)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	var userName string

	row := stmt.QueryRow(email)

	err = row.Scan(&userName)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (p *psqlUser) PsqlCheckTwoStepsVerificationIsTrue(email string) (bool, error) {
	stmt, err := p.DB.Prepare(sqlCheckUserTsvIsTrue)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)

	var ok bool
	err = row.Scan(&ok)
	if err != nil || !ok {
		return false, err
	}

	return true, nil
}

func (p *psqlUser) PsqlChangeUserName(newUserName, email string) error {
	stmt, err := p.DB.Prepare(sqlChangeUserName)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newUserName, email)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) psqlChangeUserFirstName(newFirstName, email string) error {
	stmt, err := p.DB.Prepare(sqlChangeUserFirstName)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newFirstName, email)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) psqlChangeUserLastName(newLastName, email string) error {
	stmt, err := p.DB.Prepare(sqlChangeUserLastName)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newLastName, email)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) psqlChangeUserEmail(newEmail, actualEmail, password string) error {
	stmt, err := p.DB.Prepare(sqlChangeUserEmail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newEmail, actualEmail, password)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) psqlChangeUserPassword(newPassword, email string) error {
	stmt, err := p.DB.Prepare(sqlChangeUserPassword)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPassword, email)
	if err != nil {
		return err
	}
	return nil
}
