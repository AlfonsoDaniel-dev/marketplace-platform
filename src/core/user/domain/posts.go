package user_domain

import (
	"github.com/google/uuid"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"time"
)

func (u *UserDomain) CreatePostO(email, postName string, form UserDTO.CreatePostDTO) error {

	userId, err := u.OutputInterface.PsqlGetUserIdByEmail(email)
	if err != nil {
		return err
	}

	userRepository, err := u.OutputInterface.PsqlGetUserRepositoryPath(userId)
	if err != nil {
		return err
	}

	postsDir, err := u.OutputInterface.PsqlGetUserPostsDirectory(userId)
	if err != nil {
		return err
	}

	postData, err := u.UploadsInterface.NewPost(userId, postsDir, userRepository, postName)
	if err != nil {
		return err
	}

	post := models.CreatePost{
		Id:          uuid.New(),
		CreatorId:   userId,
		ContentPath: postData.Path,
		Title:       form.Title,
		Content:     form.Body,
		CreatedAt:   time.Now().Unix(),
	}

	if postsDir == "" {
		post.UserPostsDirectory = postData.UserPostsDir

		err = u.OutputInterface.PsqlInsertUserPostsDirectory(userId, postData.UserPostsDir)
		if err != nil {
			return err
		}
	}

	err = u.OutputInterface.PsqlCreatePost(post)
	if err != nil {
		return err
	}

	return nil
}
