package user_secret

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"streamerEventViewer/pkg/models"
)

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, secret models.UserSecret) error {
	query := `
INSERT INTO user_secrets(
	id, 
	user_id,
	scopes,
	auth_token)
VALUES ($1, $2, $3, $4);
`
	_, err := r.db.ExecContext(ctx, query, secret.ID, secret.UserID, pq.Array(secret.Scopes), secret.AuthToken)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) (models.UserSecret, error) {
	query := `
SELECT id, scopes, auth_token FROM user_secrets WHERE user_id = $1;
`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return models.UserSecret{}, err
	}

	rows.Next()

	var (
		id        string
		scopes    []string
		authToken string
	)

	err = rows.Scan(&id, pq.Array(&scopes), &authToken)
	if err != nil {
		return models.UserSecret{}, err
	}

	secret := models.UserSecret{
		ID:        id,
		UserID:    userID,
		Scopes:    scopes,
		AuthToken: authToken,
	}

	return secret, nil
}
