package clip

import (
	"context"
	"database/sql"
	"streamerEventViewer/pkg/models"
)

func New(db *sql.DB) models.ClipRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, clip models.Clip) error {
	query := `
INSERT INTO clips(
	id, 
	user_id,
    streamer_id,
    external_id,
    edit_url
VALUES ($1, $2, $3, $4, $5);
`
	_, err := r.db.ExecContext(ctx, query, clip.ID, clip.UserID, clip.StreamerID, clip.ExternalID, clip.EditURL)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByUserAndStreamerID(ctx context.Context, userID, streamerID string) ([]models.Clip, error) {
	query := `
SELECT id, external_id, edit_url, views FROM clips WHERE user_id = $1 AND streamer_id = $2;
`
	rows, err := r.db.QueryContext(ctx, query, userID, streamerID)
	if err != nil {
		return nil, err
	}

	var result []models.Clip

	for rows.Next() {
		var id string
		var externalID string
		var editURL string
		var views int64

		err := rows.Scan(&id, &externalID, &editURL, &views)
		if err != nil {
			return nil, err
		}

		result = append(result, models.Clip{
			ID:         id,
			UserID:     userID,
			StreamerID: streamerID,
			EditURL:    editURL,
			ExternalID: externalID,
			Views:      views,
		})
	}

	return result, nil
}
