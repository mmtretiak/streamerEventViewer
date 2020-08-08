package models

import "context"

type Clip struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	StreamerID string `json:"streamer_id"`
	ExternalID string `json:"external_id"`
	EditURL    string `json:"edit_url"`
	// For simplifying store this with another clip info, in future make sense to move into separate table or into NoSQL database
	Views int64 `json:"views"`
}

type ClipRepository interface {
	Save(context.Context, Clip) error
	GetByUserAndStreamerID(ctx context.Context, userID, streamerID string) ([]Clip, error)
}
