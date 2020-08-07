package models

import "context"

type UsersToStreamerRepository interface {
	Save(context.Context, UserToStreamer) error
	IsExist(context.Context, UserToStreamer) (bool, error)
}

type UserToStreamer struct {
	UserID     string `json:"user_id"`
	StreamerID string `json:"streamer_id"`
}
