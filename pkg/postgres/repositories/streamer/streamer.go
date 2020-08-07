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
VALUES ($1, $2);
`
	_, err := r.db.ExecContext(ctx, query, streamer.ID, streamer.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id string) (models.Streamer, error) {
	query := `
SELECT name FROM streamers WHERE id = $1;
`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return models.Streamer{}, err
	}

	rows.Next()

	var name string

	err = rows.Scan(&name)
	if err != nil {
		return models.Streamer{}, err
	}

	streamer := models.Streamer{
		ID:   id,
		Name: name,
	}

	return streamer, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (models.Streamer, error) {
	query := `
SELECT id FROM streamers WHERE name = $1;
`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return models.Streamer{}, err
	}

	rows.Next()

	var id string

	err = rows.Scan(&id)
	if err != nil {
		return models.Streamer{}, err
	}

	streamer := models.Streamer{
		ID:   id,
		Name: name,
	}

	return streamer, nil
}
