package bgs

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Follow represents a follow relationship
type Follow struct {
	Follower  string    `json:"follower"`
	Following string    `json:"following"`
	CreatedAt time.Time `json:"createdAt"`
}

// FollowRepository handles follow data access
type FollowRepository struct {
	db *pgxpool.Pool
}

// NewFollowRepository creates a new follow repository
func NewFollowRepository(db *pgxpool.Pool) *FollowRepository {
	return &FollowRepository{db: db}
}

// CreateFollow creates a new follow relationship
func (r *FollowRepository) CreateFollow(ctx context.Context, follow *Follow) error {
	query := `
		INSERT INTO follows (follower, following, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (follower, following) DO NOTHING
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		follow.Follower,
		follow.Following,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	follow.CreatedAt = now

	return nil
}

// DeleteFollow deletes a follow relationship
func (r *FollowRepository) DeleteFollow(ctx context.Context, follower, following string) error {
	query := `
		DELETE FROM follows
		WHERE follower = $1 AND following = $2
	`
	_, err := r.db.Exec(ctx, query, follower, following)
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}
	return nil
}

// GetFollowers gets the followers of a DID
func (r *FollowRepository) GetFollowers(ctx context.Context, did string, limit, offset int) ([]*Follow, error) {
	query := `
		SELECT follower, following, created_at
		FROM follows
		WHERE following = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, did, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		var follow Follow
		err := rows.Scan(
			&follow.Follower,
			&follow.Following,
			&follow.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follow: %w", err)
		}
		follows = append(follows, &follow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating follows: %w", err)
	}
	return follows, nil
}

// GetFollowing gets the DIDs a DID is following
func (r *FollowRepository) GetFollowing(ctx context.Context, did string, limit, offset int) ([]*Follow, error) {
	query := `
		SELECT follower, following, created_at
		FROM follows
		WHERE follower = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, did, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		var follow Follow
		err := rows.Scan(
			&follow.Follower,
			&follow.Following,
			&follow.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follow: %w", err)
		}
		follows = append(follows, &follow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating follows: %w", err)
	}
	return follows, nil
}

// IsFollowing checks if a DID is following another DID
func (r *FollowRepository) IsFollowing(ctx context.Context, follower, following string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM follows
			WHERE follower = $1 AND following = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, follower, following).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if following: %w", err)
	}
	return exists, nil
}

// GetFollowersCount gets the number of followers of a DID
func (r *FollowRepository) GetFollowersCount(ctx context.Context, did string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM follows
		WHERE following = $1
	`
	var count int
	err := r.db.QueryRow(ctx, query, did).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get followers count: %w", err)
	}
	return count, nil
}

// GetFollowingCount gets the number of DIDs a DID is following
func (r *FollowRepository) GetFollowingCount(ctx context.Context, did string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM follows
		WHERE follower = $1
	`
	var count int
	err := r.db.QueryRow(ctx, query, did).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get following count: %w", err)
	}
	return count, nil
}
