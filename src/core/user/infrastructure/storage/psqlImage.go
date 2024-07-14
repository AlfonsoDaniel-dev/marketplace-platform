package Userstorage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
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

func (p *psqlUser) PsqlGetUserRepositoryPath(userId uuid.UUID) (string, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return "", err
	}

	res, err := db.RunQuery(tx, sqlGetUserRepositoryPath, userId)
	if err != nil {
		return "", err
	}

	path, err := db.ParseAnyToString(res[0])
	if err != nil {
		return "", err
	}

	return path, nil
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

func (p *psqlUser) PsqlGetUserProfilePictureData(email string) (models.ImageData, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return models.ImageData{}, err
	}

	res, err := db.RunQuery(tx, sqlGetUserProfilePictureData, email)
	if err != nil {
		return models.ImageData{}, err
	}

	data := models.ImageData{}

	err = db.MapStructValues(res, &data)
	if err != nil {
		fmt.Println(err)
		return models.ImageData{}, err
	}

	if err := tx.Commit(); err != nil {
		errStr := fmt.Sprintf("Failed to commit transaction. Error: %v", err)
		return models.ImageData{}, errors.New(errStr)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		errStr := fmt.Sprintf("Failed to commit transaction. Error: %v", err)
		return models.ImageData{}, errors.New(errStr)
	}

	return data, nil
}

func (p *psqlUser) PsqlCreateCollection(path string, form UserDTO.DbCreateCollection) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlInsertCollectionData, form.Id, form.UserId, form.CollectionName, form.Description, path, now)
	if err != nil {
		return err
	}

	return nil
}
