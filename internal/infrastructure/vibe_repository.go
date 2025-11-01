package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/domain/vibe"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresVibeRepository implements vibe.Repository using PostgreSQL.
type PostgresVibeRepository struct {
	db *pgxpool.Pool
}

// NewPostgresVibeRepository creates a new PostgresVibeRepository.
func NewPostgresVibeRepository(db *pgxpool.Pool) *PostgresVibeRepository {
	return &PostgresVibeRepository{db: db}
}

// Create persists a new vibe and updates related counters atomically.
func (r *PostgresVibeRepository) Create(ctx context.Context, v *vibe.Vibe) (*vibe.CreateResult, error) {
	if v == nil {
		return nil, &shared.ValidationError{Field: "vibe", Message: "vibe cannot be nil"}
	}

	var err error
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to begin transaction",
			Err:     err,
		}
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var (
		insertQuery = `
			INSERT INTO vibes (grumble_id, user_id, vibe_type)
			VALUES ($1, $2, $3)
			RETURNING vibe_id, voted_at
		`
		vibeID  int
		votedAt time.Time
	)

	err = tx.QueryRow(ctx, insertQuery, v.GrumbleID, v.UserID, v.Type).Scan(&vibeID, &votedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, &shared.DuplicateVibeError{
				GrumbleID: string(v.GrumbleID),
				UserID:    string(v.UserID),
			}
		}
		return nil, &shared.InternalError{
			Message: "failed to create vibe",
			Err:     err,
		}
	}

	var vibeCount int
	err = tx.QueryRow(ctx, `
		UPDATE grumbles
		SET vibe_count = vibe_count + 1
		WHERE grumble_id = $1
		RETURNING vibe_count
	`, v.GrumbleID).Scan(&vibeCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &shared.NotFoundError{
				Entity: "Grumble",
				ID:     string(v.GrumbleID),
			}
		}
		return nil, &shared.InternalError{
			Message: "failed to increment grumble vibe count",
			Err:     err,
		}
	}

	var virtuePoints int
	err = tx.QueryRow(ctx, `
		UPDATE anonymous_users
		SET virtue_points = virtue_points + 1
		WHERE user_id = $1
		RETURNING virtue_points
	`, v.UserID).Scan(&virtuePoints)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &shared.NotFoundError{
				Entity: "User",
				ID:     string(v.UserID),
			}
		}
		return nil, &shared.InternalError{
			Message: "failed to increment virtue points",
			Err:     err,
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, &shared.InternalError{
			Message: "failed to commit vibe transaction",
			Err:     err,
		}
	}
	// ensure deferred rollback does not execute after successful commit
	err = nil

	v.VibeID = shared.VibeID(vibeID)
	v.VotedAt = votedAt

	return &vibe.CreateResult{
		Vibe:         v,
		VibeCount:    vibeCount,
		VirtuePoints: virtuePoints,
	}, nil
}

// Exists checks if a user has already given a vibe to the specified grumble.
func (r *PostgresVibeRepository) Exists(ctx context.Context, grumbleID shared.GrumbleID, userID shared.UserID) (bool, error) {
	query := `
		SELECT 1
		FROM vibes
		WHERE grumble_id = $1 AND user_id = $2
		LIMIT 1
	`

	var exists int
	err := r.db.QueryRow(ctx, query, grumbleID, userID).Scan(&exists)
	if err == nil {
		return true, nil
	}
	if err == pgx.ErrNoRows {
		return false, nil
	}

	return false, &shared.InternalError{
		Message: "failed to check vibe existence",
		Err:     err,
	}
}

// CountByGrumble returns the number of vibes associated with a grumble.
func (r *PostgresVibeRepository) CountByGrumble(ctx context.Context, grumbleID shared.GrumbleID) (int, error) {
	query := "SELECT COUNT(*) FROM vibes WHERE grumble_id = $1"

	var count int
	if err := r.db.QueryRow(ctx, query, grumbleID).Scan(&count); err != nil {
		return 0, &shared.InternalError{
			Message: "failed to count vibes",
			Err:     err,
		}
	}

	return count, nil
}

// FindByUser returns vibes made by a user ordered by newest first.
func (r *PostgresVibeRepository) FindByUser(ctx context.Context, userID shared.UserID, limit int, offset int) ([]*vibe.Vibe, error) {
	baseQuery := `
		SELECT vibe_id, grumble_id, user_id, vibe_type, voted_at
		FROM vibes
		WHERE user_id = $1
		ORDER BY voted_at DESC
	`

	args := []interface{}{userID}
	argPos := 2

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	if offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to query vibes",
			Err:     err,
		}
	}
	defer rows.Close()

	var vibes []*vibe.Vibe
	for rows.Next() {
		var v vibe.Vibe
		if err := rows.Scan(&v.VibeID, &v.GrumbleID, &v.UserID, &v.Type, &v.VotedAt); err != nil {
			return nil, &shared.InternalError{
				Message: "failed to scan vibe",
				Err:     err,
			}
		}
		vibes = append(vibes, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, &shared.InternalError{
			Message: "error iterating vibes",
			Err:     err,
		}
	}

	return vibes, nil
}
