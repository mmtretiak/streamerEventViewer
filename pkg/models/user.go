package models

import "context"

type UserRepository interface {
	Save(context.Context, User) error
	GetByID(context.Context, string) (User, error)
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ThumbnailURL string `json:"thumbnail_url"`
}
