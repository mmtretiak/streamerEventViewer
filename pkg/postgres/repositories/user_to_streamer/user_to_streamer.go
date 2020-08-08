package user_to_streamer

import (
	"context"
	"database/sql"
	"streamerEventViewer/pkg/models"
)

func New(db *sql.DB) models.UsersToStreamerRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, userToStreamer models.UserToStreamer) error {
	query := `
INSERT INTO users_to_streamers(
	user_id,
	streamer_id)
VALUES ($1, $2);
`
	_, err := r.db.ExecContext(ctx, query, userToStreamer.UserID, userToStreamer.StreamerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) IsExist(ctx context.Context, userToStreamer models.UserToStreamer) (bool, error) {
	query := `
SELECT user_id, streamer_id FROM users_to_streamers WHERE user_id = $1 AND streamer_id = $2;
`

	_, err := r.db.QueryContext(ctx, query, userToStreamer.UserID, userToStreamer.StreamerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
