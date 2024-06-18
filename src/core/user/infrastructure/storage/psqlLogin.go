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
