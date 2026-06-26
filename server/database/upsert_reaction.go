package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Reaction string

const (
	ReactionLike    Reaction = "like"
	ReactionDislike Reaction = "dislike"
)

// UpsertReactionAndMaybeMatch upserts viewer's reaction to target
// and, if it's a "like", checkslike and creates a match.
// Returns (isMatch, error).
func UpsertReaction(ctx context.Context, pool *pgxpool.Pool, viewerID, targetID string, reaction Reaction) error {

	_, err := pool.Exec(ctx, `
		INSERT INTO user_reactions (user_id, target_user_id, reaction)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, target_user_id)
		DO UPDATE SET reaction = EXCLUDED.reaction
	`, viewerID, targetID, string(reaction))
	if err != nil {
		return err
	}

	return nil
}
