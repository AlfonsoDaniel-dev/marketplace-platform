package Userstorage

import (
	"github.com/google/uuid"
	"shopperia/src/External/db"
	"shopperia/src/common/models"
	"time"
)

func (p *psqlUser) PsqlCreatePost(post models.CreatePost) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = db.ExecQuery(tx, sqlCreatePost, post.Id, post.CreatorId, post.ContentPath, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) PsqlUpdatePostTitle(userId uuid.UUID, title string) error {
	tx, err := p.DB.Begin()

	now := time.Now().Unix()
	_, err = db.ExecQuery(tx, sqlUpdatePostTitle, title, now, userId)
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlUser) PsqlUpdatePostContent(user_id uuid.UUID, NewContent string) error {
	tx, err := p.DB.Begin()

	now := time.Now().Unix()

	_, err = db.ExecQuery(tx, sqlUpdatePostTitle, NewContent, now, user_id)
	if err != nil {
		return err
	}

	return nil
}
