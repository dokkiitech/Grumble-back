package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresGrumbleRepository implements grumble.Repository using PostgreSQL
type PostgresGrumbleRepository struct {
	db *pgxpool.Pool
}

// NewPostgresGrumbleRepository creates a new PostgresGrumbleRepository
func NewPostgresGrumbleRepository(db *pgxpool.Pool) *PostgresGrumbleRepository {
	return &PostgresGrumbleRepository{db: db}
}

// Create stores a new grumble
func (r *PostgresGrumbleRepository) Create(ctx context.Context, g *grumble.Grumble) error {
	query := `
		INSERT INTO grumbles (
			grumble_id, user_id, content, toxic_level, vibe_count,
			is_purified, posted_at, expires_at, is_event_grumble
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, query,
		g.GrumbleID, g.UserID, g.Content, g.ToxicLevel, g.VibeCount,
		g.IsPurified, g.PostedAt, g.ExpiresAt, g.IsEventGrumble,
	)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to create grumble",
			Err:     err,
		}
	}

	return nil
}

// FindByID retrieves a grumble by its ID
func (r *PostgresGrumbleRepository) FindByID(ctx context.Context, id shared.GrumbleID) (*grumble.Grumble, error) {
	query := `
		SELECT grumble_id, user_id, content, toxic_level, vibe_count,
		       is_purified, posted_at, expires_at, is_event_grumble
		FROM grumbles
		WHERE grumble_id = $1
	`

	var g grumble.Grumble
	err := r.db.QueryRow(ctx, query, id).Scan(
		&g.GrumbleID, &g.UserID, &g.Content, &g.ToxicLevel, &g.VibeCount,
		&g.IsPurified, &g.PostedAt, &g.ExpiresAt, &g.IsEventGrumble,
	)
	if err == pgx.ErrNoRows {
		return nil, &shared.NotFoundError{
			Entity: "Grumble",
			ID:     string(id),
		}
	}
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to find grumble",
			Err:     err,
		}
	}

	return &g, nil
}

// FindTimeline retrieves grumbles for the timeline with filtering
func (r *PostgresGrumbleRepository) FindTimeline(ctx context.Context, filter grumble.TimelineFilter) ([]*grumble.Grumble, error) {
	baseQuery := `
		SELECT grumble_id, user_id, content, toxic_level, vibe_count,
		       is_purified, posted_at, expires_at, is_event_grumble
		FROM grumbles
		WHERE 1=1
	`
	query, args, argIdx := buildTimelineFilter(baseQuery, filter)

	// Order by most recent first
	query += " ORDER BY posted_at DESC"

	// Pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to query timeline",
			Err:     err,
		}
	}
	defer rows.Close()

	var grumbles []*grumble.Grumble
	for rows.Next() {
		var g grumble.Grumble
		err := rows.Scan(
			&g.GrumbleID, &g.UserID, &g.Content, &g.ToxicLevel, &g.VibeCount,
			&g.IsPurified, &g.PostedAt, &g.ExpiresAt, &g.IsEventGrumble,
		)
		if err != nil {
			return nil, &shared.InternalError{
				Message: "failed to scan grumble",
				Err:     err,
			}
		}
		grumbles = append(grumbles, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, &shared.InternalError{
			Message: "error iterating grumbles",
			Err:     err,
		}
	}

	return grumbles, nil
}

// CountTimeline returns the total count of grumbles matching the filter
func (r *PostgresGrumbleRepository) CountTimeline(ctx context.Context, filter grumble.TimelineFilter) (int, error) {
	baseQuery := "SELECT COUNT(*) FROM grumbles WHERE 1=1"
	query, args, _ := buildTimelineFilter(baseQuery, filter)

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, &shared.InternalError{
			Message: "failed to count timeline",
			Err:     err,
		}
	}

	return count, nil
}

// Update updates an existing grumble
func (r *PostgresGrumbleRepository) Update(ctx context.Context, g *grumble.Grumble) error {
	query := `
		UPDATE grumbles
		SET content = $2, toxic_level = $3, vibe_count = $4,
		    is_purified = $5, expires_at = $6, is_event_grumble = $7
		WHERE grumble_id = $1
	`

	result, err := r.db.Exec(ctx, query,
		g.GrumbleID, g.Content, g.ToxicLevel, g.VibeCount,
		g.IsPurified, g.ExpiresAt, g.IsEventGrumble,
	)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to update grumble",
			Err:     err,
		}
	}

	if result.RowsAffected() == 0 {
		return &shared.NotFoundError{
			Entity: "Grumble",
			ID:     string(g.GrumbleID),
		}
	}

	return nil
}

// ArchiveExpired moves expired grumbles to archive table and removes them from main table
func (r *PostgresGrumbleRepository) ArchiveExpired(ctx context.Context) (int, error) {
	// トランザクション開始
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, &shared.InternalError{
			Message: "failed to start transaction",
			Err:     err,
		}
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	// 1. 期限切れの投稿をアーカイブテーブルに挿入
	insertQuery := `
		INSERT INTO grumbles_archive
			(grumble_id, user_id, content, toxic_level, vibe_count,
			 is_purified, posted_at, expires_at, is_event_grumble, archived_at)
		SELECT
			grumble_id, user_id, content, toxic_level, vibe_count,
			is_purified, posted_at, expires_at, is_event_grumble, $1
		FROM grumbles
		WHERE expires_at <= $2
	`

	_, err = tx.Exec(ctx, insertQuery, now, now)
	if err != nil {
		return 0, &shared.InternalError{
			Message: "failed to archive expired grumbles",
			Err:     err,
		}
	}

	// 2. grumblesテーブルから削除
	deleteQuery := "DELETE FROM grumbles WHERE expires_at <= $1"
	result, err := tx.Exec(ctx, deleteQuery, now)
	if err != nil {
		return 0, &shared.InternalError{
			Message: "failed to delete expired grumbles",
			Err:     err,
		}
	}

	// トランザクションコミット
	if err := tx.Commit(ctx); err != nil {
		return 0, &shared.InternalError{
			Message: "failed to commit transaction",
			Err:     err,
		}
	}

	return int(result.RowsAffected()), nil
}

// FindPurificationCandidates finds grumbles that meet purification threshold
func (r *PostgresGrumbleRepository) FindPurificationCandidates(ctx context.Context, threshold int) ([]*grumble.Grumble, error) {
	query := `
		SELECT grumble_id, user_id, content, toxic_level, vibe_count,
		       is_purified, posted_at, expires_at, is_event_grumble
		FROM grumbles
		WHERE is_purified = FALSE AND vibe_count >= $1
	`

	rows, err := r.db.Query(ctx, query, threshold)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to find purification candidates",
			Err:     err,
		}
	}
	defer rows.Close()

	var grumbles []*grumble.Grumble
	for rows.Next() {
		var g grumble.Grumble
		err := rows.Scan(
			&g.GrumbleID, &g.UserID, &g.Content, &g.ToxicLevel, &g.VibeCount,
			&g.IsPurified, &g.PostedAt, &g.ExpiresAt, &g.IsEventGrumble,
		)
		if err != nil {
			return nil, &shared.InternalError{
				Message: "failed to scan grumble",
				Err:     err,
			}
		}
		grumbles = append(grumbles, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, &shared.InternalError{
			Message: "error iterating purification candidates",
			Err:     err,
		}
	}

	return grumbles, nil
}

// IncrementVibeCount atomically increments the vibe count for a grumble
func (r *PostgresGrumbleRepository) IncrementVibeCount(ctx context.Context, id shared.GrumbleID) error {
	query := "UPDATE grumbles SET vibe_count = vibe_count + 1 WHERE grumble_id = $1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return &shared.InternalError{
			Message: "failed to increment vibe count",
			Err:     err,
		}
	}

	if result.RowsAffected() == 0 {
		return &shared.NotFoundError{
			Entity: "Grumble",
			ID:     string(id),
		}
	}

	return nil
}

// FindArchivedTimeline retrieves grumbles from archive table for a specific date
func (r *PostgresGrumbleRepository) FindArchivedTimeline(
	ctx context.Context,
	filter grumble.TimelineFilter,
	targetDate time.Time,
) ([]*grumble.Grumble, error) {
	// 対象日の00:00〜23:59を取得
	dayStart, dayEnd := shared.GetDayBounds(targetDate)

	baseQuery := `
		SELECT grumble_id, user_id, content, toxic_level, vibe_count,
		       is_purified, posted_at, expires_at, is_event_grumble
		FROM grumbles_archive
		WHERE posted_at >= $1 AND posted_at <= $2
	`

	args := []interface{}{dayStart, dayEnd}
	argIdx := 3

	// フィルタ条件を追加
	query := baseQuery

	if filter.ToxicLevelMin != nil {
		query += fmt.Sprintf(" AND toxic_level >= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMin)
		argIdx++
	}

	if filter.ToxicLevelMax != nil {
		query += fmt.Sprintf(" AND toxic_level <= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMax)
		argIdx++
	}

	if filter.ExcludePurified {
		query += " AND is_purified = FALSE"
	}

	// ソートとページネーション
	query += " ORDER BY posted_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to query archived timeline",
			Err:     err,
		}
	}
	defer rows.Close()

	var grumbles []*grumble.Grumble
	for rows.Next() {
		var g grumble.Grumble
		err := rows.Scan(
			&g.GrumbleID, &g.UserID, &g.Content, &g.ToxicLevel, &g.VibeCount,
			&g.IsPurified, &g.PostedAt, &g.ExpiresAt, &g.IsEventGrumble,
		)
		if err != nil {
			return nil, &shared.InternalError{
				Message: "failed to scan archived grumble",
				Err:     err,
			}
		}
		grumbles = append(grumbles, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, &shared.InternalError{
			Message: "error iterating archived grumbles",
			Err:     err,
		}
	}

	return grumbles, nil
}

// CountArchivedTimeline counts archived grumbles for a specific date
func (r *PostgresGrumbleRepository) CountArchivedTimeline(
	ctx context.Context,
	filter grumble.TimelineFilter,
	targetDate time.Time,
) (int, error) {
	dayStart, dayEnd := shared.GetDayBounds(targetDate)

	query := `
		SELECT COUNT(*)
		FROM grumbles_archive
		WHERE posted_at >= $1 AND posted_at <= $2
	`

	args := []interface{}{dayStart, dayEnd}
	argIdx := 3

	if filter.ToxicLevelMin != nil {
		query += fmt.Sprintf(" AND toxic_level >= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMin)
		argIdx++
	}

	if filter.ToxicLevelMax != nil {
		query += fmt.Sprintf(" AND toxic_level <= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMax)
		argIdx++
	}

	if filter.ExcludePurified {
		query += " AND is_purified = FALSE"
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, &shared.InternalError{
			Message: "failed to count archived timeline",
			Err:     err,
		}
	}

	return count, nil
}

func buildTimelineFilter(base string, filter grumble.TimelineFilter) (string, []interface{}, int) {
	query := base
	args := []interface{}{}
	argIdx := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, string(*filter.UserID))
		argIdx++
	}

	if filter.ExcludePurified {
		query += " AND is_purified = FALSE"
	}

	if filter.ExcludeExpired {
		query += fmt.Sprintf(" AND expires_at > $%d", argIdx)
		args = append(args, time.Now())
		argIdx++
	}

	if filter.ToxicLevelMin != nil {
		query += fmt.Sprintf(" AND toxic_level >= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMin)
		argIdx++
	}

	if filter.ToxicLevelMax != nil {
		query += fmt.Sprintf(" AND toxic_level <= $%d", argIdx)
		args = append(args, *filter.ToxicLevelMax)
		argIdx++
	}

	return query, args, argIdx
}
