package Userstorage

import (
	"fmt"
	"shopperia/src/db"
)

func (p *psqlUser) PsqlGetHashPassword(email string) ([]byte, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	res, err := db.RunQuery(tx, sqlGetHashedPasswordFromEmail, email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	hashPassword, err := db.ParseAnyToString(res[0])
	if err != nil {
		return nil, err
	}

	return []byte(hashPassword), nil
}

func (p *psqlUser) PsqlVerifyEmailExists(email string) (string, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return "", err
	}

	res, err := db.RunQuery(tx, sqlLoginVerifyEmailExists, email)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	fmt.Println(res[0])
	existingEmail, err := db.ParseAnyToString(res[0])
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return existingEmail, nil
}

func (p *psqlUser) PsqlInsertTsvCode(email, code string) (string, error) {
	stmt, err := p.DB.Prepare(sqlInsertAccesToken)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	row := stmt.QueryRow(code, email)

	var accesToken string

	err = row.Scan(&accesToken)
	if err != nil {
		return "", err
	}

	return accesToken, nil
}

func (p *psqlUser) PsqlCleanAccessToken(email string) error {
	stmt, err := p.DB.Prepare(sqlCleanAccessToken)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) PsqlGetAccessToken(email string) (string, error) {
	stmt, err := p.DB.Prepare(sqlGetAccessToken)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)

	var accessToken string

	err = row.Scan(&accessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
