package Userstorage

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	user_model "shopperia/src/common/models"
	"shopperia/src/core/helpers"
	"shopperia/src/db"
	"time"
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
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	nullTime := helpers.IntToNull(user.UpdatedAt)

	nullBiography := helpers.StringToNull(user.Biography)

	if _, err := db.ExecQuery(tx, sqlCreateUser, user.Id, user.FirstName, user.LastName, user.UserName, nullBiography, user.Age, user.Email, user.Password, user.TwoStepsVerfication, user.CreatedAt, nullTime); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("User created successfully")
	return nil
}

func (p *psqlUser) PsqlGetUserIdByEmail(email string) (uuid.UUID, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	res, err := db.RunQuery(tx, sqlGetUserIdByEmail, email)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	id, err := db.ParseAnyToUUID(res)

	return id, nil
}

func (p *psqlUser) PsqlInsertAddressData(address user_model.Address) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = db.ExecQuery(tx, sqlInsertUserAddress)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (p *psqlUser) PsqlGetUserNameByEmail(email string) (string, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return "", err
	}

	res, err := db.RunQuery(tx, sqlGetUserNameByEmail, email)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	userName, err := db.ParseAnyToString(res[0])
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return userName, nil
}

func (p *psqlUser) PsqlCheckTwoStepsVerificationIsTrue(email string) (bool, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return false, err
	}

	res, err := db.RunQuery(tx, sqlCheckUserTsvIsTrue, email)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	status, err := db.ParseAnyToBool(res[0])
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()

	if !status {
		return false, nil
	}

	return true, nil
}

func (p *psqlUser) PsqlChangeUserName(newUserName, email string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlChangeUserName, newUserName, now, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (p *psqlUser) PsqlChangeUserFirstName(newFirstName, email string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlChangeUserFirstName, newFirstName, now, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *psqlUser) PsqlChangeUserLastName(newLastName, email string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	now := time.Now().Unix()

	_, err = db.ExecQuery(tx, sqlChangeUserLastName, newLastName, now, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *psqlUser) PsqlChangeUserEmail(newEmail, userId string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if _, err := db.ExecQuery(tx, sqlChangeUserEmail, newEmail, now, userId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *psqlUser) PsqlChangeUserPassword(newPassword, email, oldPassword string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := db.ExecQuery(tx, sqlChangeUserPassword, newPassword, email, oldPassword); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *psqlUser) PsqlChangeUserTsvStatus(email string, value bool) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := db.ExecQuery(tx, sqlChangeTsvStatus, value, email); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *psqlUser) PsqlInsertProfilePictureData(fileName, filePath, )
