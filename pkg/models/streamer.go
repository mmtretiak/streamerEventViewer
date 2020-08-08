package models

import "context"

type StreamerRepository interface {
	Save(context.Context, Streamer) error
	GetByID(context.Context, string) (Streamer, error)
	GetByName(context.Context, string) (Streamer, error)
	GetByUserID(context.Context, string) ([]Streamer, error)
}

type Streamer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
}
