package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
var ErrProfileNotFound = errors.New("profile not found")
// UpdateParentProfileDynamic updates with prebuilt SET clauses and args,
// and returns a map keyed by the RETURNING column names.
func UpdateProfileDynamic(
	ctx context.Context,
	pool *pgxpool.Pool,
	sets []string,
	args []any,
	returningCols []string,
	table string,
) (map[string]any, error) {

	// return user_id as text with a stable key
	ret := make([]string, 0, len(returningCols)+1)
	ret = append(ret, `user_id::text AS user_id`)
	ret = append(ret, returningCols...)

	q := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE user_id = $1
		RETURNING %s`,
		table,
		strings.Join(sets, ", "),
		strings.Join(ret, ", "),
	)

	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("update %s: %w", table, err)
	}
	defer rows.Close()

	m, err := pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("scan %s: %w", table, err)
	}
	return m, nil
}

