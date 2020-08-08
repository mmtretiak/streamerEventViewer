package streamer

import (
	"context"
	"database/sql"
	"streamerEventViewer/pkg/models"
)

func New(db *sql.DB) models.StreamerRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, streamer models.Streamer) error {
	query := `
INSERT INTO streamers(
	id, 
	name,
    external_id)
VALUES ($1, $2, $3);
`
	_, err := r.db.ExecContext(ctx, query, streamer.ID, streamer.Name, streamer.ExternalID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id string) (models.Streamer, error) {
	query := `
SELECT external_id, name FROM streamers WHERE id = $1;
`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return models.Streamer{}, err
	}

	rows.Next()

	var externalID string
	var name string

	err = rows.Scan(&externalID, &name)
	if err != nil {
		return models.Streamer{}, err
	}

	streamer := models.Streamer{
		ID:         id,
		Name:       name,
		ExternalID: externalID,
	}

	return streamer, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (models.Streamer, error) {
	query := `
SELECT id, external_id FROM streamers WHERE name = $1;
`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return models.Streamer{}, err
	}

	rows.Next()

	var id string
	var externalID string

	err = rows.Scan(&id, &externalID)
	if err != nil {
		return models.Streamer{}, err
	}

	streamer := models.Streamer{
		ID:         id,
		Name:       name,
		ExternalID: externalID,
	}

	return streamer, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) ([]models.Streamer, error) {
	query := `
SELECT streamers.id, streamers.name, streamers.external_id FROM streamers LEFT JOIN users_to_streamers on streamers.id = users_to_streamers.streamer_id WHERE user_id = $1;
`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var result []models.Streamer

	for rows.Next() {
		var id string
		var name string
		var externalID string

		err := rows.Scan(&id, &name, &externalID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.Streamer{
			ID:         id,
			Name:       name,
			ExternalID: externalID,
		})
	}

	return result, nil
}
