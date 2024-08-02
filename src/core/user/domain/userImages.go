package user_domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (u *UserDomain) GetProfilePicture(email string) (models.GetImage, error) {
	if email == "" {
		return models.GetImage{}, errors.New("no email provide")
	}

	ProfilePicData, err := u.OutputInterface.PsqlGetUserProfilePictureData(email)
	if err != nil {
		return models.GetImage{}, err
	}

	image, err := u.UploadsInterface.GetProfilePicture(ProfilePicData.UserMediaRepository, ProfilePicData.FileName, ProfilePicData.FileExtension, ProfilePicData.UserId)
	if err != nil {
		fmt.Println(err)
		return models.GetImage{}, err
	}

	return image, nil
}

func (u *UserDomain) CreateCollection(form UserDTO.CreateCollection) error {
	if form.Email == "" || form.CollectionName == "" {
		return errors.New("no collection name")
	}

	userId, err := u.OutputInterface.PsqlGetUserIdByEmail(form.Email)
	if err != nil {
		return err
	}

	userName, err := u.OutputInterface.PsqlGetUserNameByEmail(form.Email)
	if err != nil {
		return err
	}

	repositoryPath, err := u.OutputInterface.PsqlGetUserRepositoryPath(userId)
	if err != nil {
		return err
	}

	collectionData, err := u.UploadsInterface.CreateCollection(userId, form.CollectionName, repositoryPath)

	collectionForm := UserDTO.DbCreateCollection{
		Id:                 uuid.New(),
		UserId:             userId,
		UserName:           userName,
		CollectionName:     form.CollectionName,
		Description:        form.Description,
		UserRepositoryPath: repositoryPath,
	}

	err = u.OutputInterface.PsqlCreateCollection(collectionData.CollectionPath, collectionForm)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDomain) UploadNewImage(imageform models.UploadImageForm) error {
	if imageform.FileName == "" || imageform.FileExtension == "" || imageform.UserEmail == "" {
		return errors.New("no needed fields provide")
	}

	userId, err := u.GetUserIdByEmail(imageform.UserEmail)
	if err != nil {
		return err
	}

	imageform.UserID = userId
	repositoryPath, err := u.OutputInterface.PsqlGetUserRepositoryPath(imageform.UserID)

	imageData, err := u.UploadsInterface.UploadMedia(repositoryPath, imageform)
	if err != nil {
		return err
	}

	err = u.OutputInterface.PsqlInsertImageData(imageData.Image_id, imageData.UserId, imageData.UserMediaRepository, imageData.FileName, imageform.FileExtension, imageData.ImagePath)
	if err != nil {
		return err
	}

	return nil
}
