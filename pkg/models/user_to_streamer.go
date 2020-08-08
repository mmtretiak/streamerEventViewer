package models

import "context"

type UsersToStreamerRepository interface {
	Save(context.Context, UserToStreamer) error
	IsExist(context.Context, UserToStreamer) (bool, error)
}

type UserToStreamer struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	StreamerID string `json:"streamer_id"`
}
