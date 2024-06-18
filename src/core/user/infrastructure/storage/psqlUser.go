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

	nullTime := helpers.TimeToNull(user.UpdatedAt)

	nullBiography := helpers.StringToNull(user.Biography)

	_, err = stmt.Exec(user.Id, user.FirstName, user.LastName, user.UserName, nullBiography, user.Age, user.Email, user.Password, user.CreatedAt, nullTime)
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

func (p *psqlUser) PsqlInsertAddressData(userId uuid.UUID, address user_model.Address) error {
	stmt, err := p.DB.Prepare(sqlInsertAddressOnAddressTable)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(address.ID, userId, address.Street, address.City, address.State, address.PostalCode, address.Country, address.CreatedAt)
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
