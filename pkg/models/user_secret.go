package models

import "context"

type UserSecretRepository interface {
	Save(context.Context, UserSecret) error
	GetByUserID(context.Context, string) (UserSecret, error)
}

type UserSecret struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	Scopes    []string `json:"scopes"`
	AuthToken string   `json:"auth_token"`
}
