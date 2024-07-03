package Userstorage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/db"
	"time"
)

func (p *psqlUser) PsqlUploadImageData(id, user_id uuid.UUID, userRepository, fileName, fileExtension, filePath string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	createdAt := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlInsertImageData, id, user_id, userRepository, fileName, fileExtension, filePath, createdAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		errStr := fmt.Sprintf("Failed to commit transaction. Error: %v", err)
		return errors.New(errStr)
	}

	return nil
}

func (p *psqlUser) PsqlInsert