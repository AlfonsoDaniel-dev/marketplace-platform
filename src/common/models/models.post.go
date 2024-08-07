package models

import "github.com/google/uuid"

type Post struct {
	ID              uuid.UUID `json:"id"`
	CreatorId       uuid.UUID `json:"creator_id"`
	CreatorUserName string    `json:"creator_user_name"`
	ContentPath     string
	Title           string `json:"title"`
	Body            string `json:"body"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type PostData struct {
	ID             uuid.UUID
	UserRepository string
	UserPostsDir   string
	Path           string
}

type CreatePost struct {
	Id                 uuid.UUID
	CreatorId          uuid.UUID
	UserPostsDirectory string
	ContentPath        string
	Title              string
	Content            string
	CreatedAt          int64
	UpdatedAt          int64
}
