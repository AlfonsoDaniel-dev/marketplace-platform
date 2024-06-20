package Userstorage

func (p *psqlUser) PsqlGetHashPassword(email string) ([]byte, error) {
	stmt, err := p.DB.Prepare(sqlGetHashedPasswordFromEmail)
	if err != nil {
		return []byte(""), err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)
	if err != nil {
		return []byte(""), err
	}

	var passwordFromDB string

	err = row.Scan(&passwordFromDB)

	hashPassword := []byte(passwordFromDB)

	return hashPassword, nil
}

func (p *psqlUser) PsqlVerifyEmailExists(email string) (string, error) {
	stmt, err := p.DB.Prepare(sqlLoginVerifyEmailExists)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)

	var existingEmail string

	err = row.Scan(&existingEmail)
	if err != nil {
		return "", err
	}

	return email, nil
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
