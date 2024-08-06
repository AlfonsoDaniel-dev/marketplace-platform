package UserDTO

import "time"

type CreatePostDTO struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type GetPostDTO struct {
	ID           int       `json:"id"`
	CreatorEmail string    `json:"creator_email"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdatePostDTO struct {
	NewTitle string `json:"new_title"`
	NewBody  string `json:"new_body"`
}
