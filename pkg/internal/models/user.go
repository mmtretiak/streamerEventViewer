package models

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ThumbnailURL string `json:"thumbnail_url"`
}
