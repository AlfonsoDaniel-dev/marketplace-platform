package Userstorage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/db"
	"time"
)

func (p *psqlUser) PsqlInsertRepositoryPathOnUser(userId uuid.UUID, repositoryPath string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = db.ExecQuery(tx, sqlInsertRepositoryPathOnUser, repositoryPath, userId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		errStr := fmt.Sprintf("falied to commit transaction: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (p *psqlUser) PsqlInsertImageData(imageID, userId uuid.UUID, userRepositoryPath, fileName, fileExtension, filePath string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	createdAt := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlInsertImageData, imageID, userId, userRepositoryPath, fileName, fileExtension, filePath, createdAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		errStr := fmt.Sprintf("Failed to commit transaction. Error: %v", err)
		return errors.New(errStr)
	}

	return nil
}
