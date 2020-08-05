package user

import (
	"context"
	"database/sql"
	"streamerEventViewer/pkg/models"
)

func New(db *sql.DB) models.UserRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, user models.User) error {
	query := `
INSERT INTO users(
	id, 
	name,
	email,
	thumbnail_url)
VALUES ($1, $2, $3, $4);
`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.ThumbnailURL)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id string) (models.User, error) {
	query := `
SELECT name, email, thumbnail_url FROM users WHERE id = $1;
`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return models.User{}, err
	}

	rows.Next()

	var (
		name         string
		email        string
		thumbnailUrl string
	)

	err = rows.Scan(&name, &email, &thumbnailUrl)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		ThumbnailURL: thumbnailUrl,
	}

	return user, nil
}
