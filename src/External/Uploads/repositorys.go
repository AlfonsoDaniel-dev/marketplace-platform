package Uploads

import (
	"errors"
	"github.com/google/uuid"
	"os"
)

func (US *UploadService) MakeNewMediaRepositoryForUser(userId uuid.UUID) (string, error) {
	mediaPath := "/user" + "_" + userId.String()

	mediaPath, err := US.MakeNewDirectory("", mediaPath)
	if err != nil {
		return mediaPath, err
	}

	return mediaPath, nil
}

func (US *UploadService) CheckUserHasAMediaRepository(userId uuid.UUID) bool {
	repoPath := getEntryPoint() + US.OriginPath + "/" + "user" + "_" + userId.String()
	_, err := os.Stat(repoPath)
	if os.IsExist(err) {
		return true
	}

	return false
}

func (US *UploadService) GetUserRepositoryPath(userId uuid.UUID) (string, error) {
	ok := US.CheckUserHasAMediaRepository(userId)
	if !ok {
		return "", errors.New("user doesnt have a media repository does not exist")
	}

	path := US.OriginPath + "/" + "user" + "_" + userId.String()

	return path, nil
}
