package Userstorage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	db2 "shopperia/src/External/db"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"time"
)

func (p *psqlUser) PsqlInsertRepositoryPathOnUser(userId uuid.UUID, repositoryPath string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = db2.ExecQuery(tx, sqlInsertRepositoryPathOnUser, repositoryPath, userId)
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

	res, err := db2.RunQuery(p.DB, sqlGetUserRepositoryPath, userId)
	if err != nil {
		return "", err
	}

	path, err := db2.ParseAnyToString(res[0])
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
	_, err = db2.ExecQuery(tx, sqlInsertImageData, imageID, userId, userRepositoryPath, fileName, fileExtension, filePath, createdAt)
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
	res, err := db2.RunQuery(p.DB, sqlGetUserProfilePictureData, email)
	if err != nil {
		return models.ImageData{}, err
	}

	data := models.ImageData{}

	err = db2.MapStructValues(res, &data)
	if err != nil {
		fmt.Println(err)
		return models.ImageData{}, err
	}

	return data, nil
}

func (p *psqlUser) PsqlCreateCollection(path string, form UserDTO.DbCreateCollection) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	_, err = db2.ExecQuery(tx, sqlInsertCollectionData, form.Id, form.UserId, form.CollectionName, form.Description, path, now)
	if err != nil {
		return err
	}

	return nil
}
