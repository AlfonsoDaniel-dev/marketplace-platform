package user_domain

import (
	"errors"
	"github.com/google/uuid"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

/*
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
} */

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
