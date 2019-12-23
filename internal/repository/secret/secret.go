package secret

import (
	"context"
	"time"

	secretentity "github.com/albertwidi/go-project-example/internal/entity/secret"
	"github.com/albertwidi/go-project-example/internal/pkg/sqldb"
	"github.com/lib/pq"
)

// Repository of secret
type Repository struct {
	db *sqldb.DB
}

// New secret
func New(db *sqldb.DB) *Repository {
	r := Repository{
		db: db,
	}
	return &r
}

// Secret of user
type Secret struct {
	ID          string           `db:"id"`
	UserID      string           `db:"user_id"`
	SecretKey   secretentity.Key `db:"secret_key"`
	SecretValue string           `db:'secret_value"`
	CreatedAt   time.Time        `db:"created_at"`
	CreatedBy   int64            `db:"created_by"`
	UpdatedAt   pq.NullTime      `db:"updated_at"`
	UpdatedBy   int64            `db:"updated_by"`
	IsTest      bool             `db:"is_test"`
}

const (
	createSecretQuery = `
INSERT INTO user_secrets(
	user_id,
	secret_key,
	secret_value,
	created_at,
	created_by,
	is_test
)
VALUES(?, ?, ?, ?, ?, ?)
`
)

// Create secret based on secret data
func (r Repository) Create(ctx context.Context, secret secretentity.Secret) error {
	q := r.db.Rebind(createSecretQuery)

	createdAt := time.Now()
	_, err := r.db.ExecContext(ctx, q, secret.UserID, secret.SecretKey, secret.SecretValue, secret.CreatedAt, createdAt)
	if err != nil {
		return err
	}

	return nil
}

const (
	getSecretQuery = `
		SELECT id, user_id, secret_key, secret_value, created_at, created_by, updated_at, updated_by, is_test
		FROM user_secrets 
		WHERE user_id = $1
			AND secret_key = $2
		`
)

// GetSecret for getting a secret
func (r Repository) GetSecret(ctx context.Context, userID int64, secretKey secretentity.Key) (*secretentity.Secret, error) {
	s := Secret{}
	err := r.db.Get(&s, getSecretQuery, userID, secretKey)
	if err != nil {
		return nil, err
	}

	return &secretentity.Secret{
		ID:        s.ID,
		UserID:    s.UserID,
		SecretKey: s.SecretKey,
		CreatedAt: s.CreatedAt,
		CreatedBy: s.CreatedBy,
		UpdatedAt: s.UpdatedAt.Time,
		UpdatedBy: s.UpdatedBy,
		IsTest:    s.IsTest,
	}, nil
}

const (
	updateSecretQuery = `
UPDATE user_secrets
SET secret_value = :secret_value,
updated_at = :updated_at,
updated_by = :updated_by
WHERE id = :id
AND secret_key = :secret_key
`
)

// UpdateSecret for updating current secret
func (r Repository) UpdateSecret(ctx context.Context, secret secretentity.Secret) error {
	_, err := r.db.NamedExecContext(ctx, updateSecretQuery, secret)
	if err != nil {
		return err
	}

	return nil
}
