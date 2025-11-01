package infrastructure

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository implements user.Repository using PostgreSQL
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create stores a new anonymous user
func (r *PostgresUserRepository) Create(ctx context.Context, u *user.AnonymousUser) error {
	query := `
		INSERT INTO anonymous_users (user_id, virtue_points, created_at, profile_title)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, u.UserID, u.VirtuePoints, u.CreatedAt, u.ProfileTitle)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to create user",
			Err:     err,
		}
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id shared.UserID) (*user.AnonymousUser, error) {
	query := `
		SELECT user_id, virtue_points, created_at, profile_title
		FROM anonymous_users
		WHERE user_id = $1
	`

	var u user.AnonymousUser
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.UserID, &u.VirtuePoints, &u.CreatedAt, &u.ProfileTitle,
	)
	if err == pgx.ErrNoRows {
		return nil, &shared.NotFoundError{
			Entity: "User",
			ID:     string(id),
		}
	}
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to find user",
			Err:     err,
		}
	}

	return &u, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, u *user.AnonymousUser) error {
	query := `
		UPDATE anonymous_users
		SET virtue_points = $2, profile_title = $3
		WHERE user_id = $1
	`

	result, err := r.db.Exec(ctx, query, u.UserID, u.VirtuePoints, u.ProfileTitle)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to update user",
			Err:     err,
		}
	}

	if result.RowsAffected() == 0 {
		return &shared.NotFoundError{
			Entity: "User",
			ID:     string(u.UserID),
		}
	}

	return nil
}

// FindTopByVirtuePoints retrieves top users by virtue points for rankings
func (r *PostgresUserRepository) FindTopByVirtuePoints(ctx context.Context, limit int) ([]*user.AnonymousUser, error) {
	query := `
		SELECT user_id, virtue_points, created_at, profile_title
		FROM anonymous_users
		ORDER BY virtue_points DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to query top users",
			Err:     err,
		}
	}
	defer rows.Close()

	var users []*user.AnonymousUser
	for rows.Next() {
		var u user.AnonymousUser
		err := rows.Scan(&u.UserID, &u.VirtuePoints, &u.CreatedAt, &u.ProfileTitle)
		if err != nil {
			return nil, &shared.InternalError{
				Message: "failed to scan user",
				Err:     err,
			}
		}
		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, &shared.InternalError{
			Message: "error iterating users",
			Err:     err,
		}
	}

	return users, nil
}

// IncrementVirtuePoints atomically increments a user's virtue points
func (r *PostgresUserRepository) IncrementVirtuePoints(ctx context.Context, id shared.UserID, points int) error {
	query := "UPDATE anonymous_users SET virtue_points = virtue_points + $2 WHERE user_id = $1"

	result, err := r.db.Exec(ctx, query, id, points)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to increment virtue points",
			Err:     err,
		}
	}

	if result.RowsAffected() == 0 {
		return &shared.NotFoundError{
			Entity: "User",
			ID:     string(id),
		}
	}

	return nil
}
