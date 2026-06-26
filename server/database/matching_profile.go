package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchingProfile struct {
	UserID            string   `db:"user_id"`
	ParentName        string   `db:"parent_name"`
	City              string   `db:"address_city"`
	Languages         []string `db:"languages"`
	PreferredDistance int      `db:"preferred_distance_km"`

	// Child information
	ChildName     string    `db:"child_name"`
	ChildBirthday time.Time `db:"child_birthday"`
	ChildGender   string    `db:"child_gender"`
	Interests     []string  `db:"interests"`
	ActivityLevel string    `db:"activity_level"`
	Limitations   []string  `db:"limitations"`
	Allergies     []string  `db:"allergies"`
	PlayStyles    []string  `db:"play_styles"`
}


// Looking for matching. COALESCE to be sure that we get something as a result
func GetMatchingProfile(ctx context.Context, pool *pgxpool.Pool, userID string) (*MatchingProfile, error) {
	const query = `
		SELECT 
			pp.user_id::text,
			COALESCE(pp.name, '') as parent_name,
			COALESCE(pp.address_city, '') as address_city,
			COALESCE(pp.languages, '{}') as languages,
			COALESCE(pp.preferred_distance_km, 0) as preferred_distance_km,
			COALESCE(c.name, '') as child_name,
			COALESCE(c.birthday, now()::date) as child_birthday,
			COALESCE(c.gender, '') as child_gender,
			COALESCE(c.interests, '{}') as interests,
			COALESCE(c.activity_level, '') as activity_level,
			COALESCE(c.limitations, '{}') as limitations,
			COALESCE(c.allergies, '{}') as allergies,
			COALESCE(c.play_styles, '{}') as play_styles
		FROM parent_profiles pp
		JOIN children c ON pp.user_id = c.user_id
		WHERE pp.user_id = $1  
		LIMIT 1` //for 1 parent 1 child

	var profile MatchingProfile
	err := pool.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.ParentName,
		&profile.City,
		&profile.Languages,
		&profile.PreferredDistance,
		&profile.ChildName,
		&profile.ChildBirthday,
		&profile.ChildGender,
		&profile.Interests,
		&profile.ActivityLevel,
		&profile.Limitations,
		&profile.Allergies,
		&profile.PlayStyles,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Profile not found
		}
		return nil, err
	}

	return &profile, nil
}



// GetPotentialMatches retrieves potential matching candidates for a user
// Excludes already connected users and dismissed recommendations
func GetPotentialMatches(ctx context.Context, pool *pgxpool.Pool, userID string) ([]MatchingProfile, error) {
	const query = `
WITH pf AS (
  -- prefilter: only candidates within viewer's radius & with filled profiles
  SELECT candidate_user_id
  FROM prefilter_candidates_postgis($1::uuid, 200, 0)
)
SELECT 
  pp.user_id::text,
  COALESCE(pp.name, '')                          AS parent_name,
  COALESCE(pp.address_city, '')                  AS address_city,
  COALESCE(pp.languages, '{}'::text[])           AS languages,
  COALESCE(pp.preferred_distance_km, 0)          AS preferred_distance_km,
  COALESCE(c.name, '')                           AS child_name,
  COALESCE(c.birthday, now()::date)              AS child_birthday,
  COALESCE(c.gender, '')                         AS child_gender,
  COALESCE(c.interests, '{}'::text[])            AS interests,
  COALESCE(c.activity_level, '')                 AS activity_level,
  COALESCE(c.limitations, '{}'::text[])          AS limitations,
  COALESCE(c.allergies, '{}'::text[])            AS allergies,
  COALESCE(c.play_styles, '{}'::text[])          AS play_styles
FROM pf
JOIN parent_profiles pp ON pp.user_id = pf.candidate_user_id
JOIN children        c  ON c.user_id  = pp.user_id
WHERE pp.user_id <> $1::uuid
  AND pp.user_id NOT IN (
    -- Exclude users already reacted with viewer
    SELECT target_user_id
    FROM user_reactions
    WHERE user_id = $1::uuid AND reaction IN ('like', 'dislike')
    UNION
    SELECT user_id
    FROM user_reactions
    WHERE target_user_id = $1::uuid AND reaction IN ('like', 'dislike')
  )
`;

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []MatchingProfile
	for rows.Next() {
		var profile MatchingProfile
		err := rows.Scan(
			&profile.UserID,
			&profile.ParentName,
			&profile.City,
			&profile.Languages,
			&profile.PreferredDistance,
			&profile.ChildName,
			&profile.ChildBirthday,
			&profile.ChildGender,
			&profile.Interests,
			&profile.ActivityLevel,
			&profile.Limitations,
			&profile.Allergies,
			&profile.PlayStyles,
		)
		if err != nil {
			continue 
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}
