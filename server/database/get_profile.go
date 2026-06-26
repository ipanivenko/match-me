package database

import (
	"context"
	"errors"
	"log"
	"matchme-server/structs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserProfile(ctx context.Context, pool *pgxpool.Pool, id string) (*structs.ParentProfile, error){
	var p structs.ParentProfile
	row := pool.QueryRow(ctx, `
        SELECT 
            user_id::text,
            COALESCE(name, '') AS name,
            COALESCE(gender, '') AS gender,
            COALESCE(about, '') AS about,
            COALESCE(languages, '{}'::text[]) AS languages,
            COALESCE(address_city, '') AS address_city,
            COALESCE(lat, 0.0)                AS lat,
            COALESCE(lon, 0.0)                AS lon,
            COALESCE(preferred_distance_km, 0) AS preferred_distance_km
        FROM parent_profiles
        WHERE user_id = $1
        LIMIT 1`, id)

    err := row.Scan(
        &p.UserID,
        &p.Name,
        &p.Gender,
        &p.About,
        &p.Languages,
        &p.AddressCity,
        &p.Lat,
        &p.Lon,     
        &p.PreferredDistance,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil //not found
        }
        return nil, err
    }

    return &p, nil
}

func GetChildProfile(ctx context.Context, pool *pgxpool.Pool, id string) (*structs.Child, error){
	var c structs.Child
	row := pool.QueryRow(ctx, `
        SELECT 
            user_id::text,
            COALESCE(name, '') AS name,
            COALESCE(birthday, '0001-01-01'::date) AS birthday,
            COALESCE(gender, '') AS gender,
            COALESCE(about_short, '') AS about_short,
            COALESCE(interests, '{}'::text[]) AS interests,
            COALESCE(activity_level, '') AS activity_level,
            COALESCE(limitations, '{}'::text[]) AS limitations,
            COALESCE(allergies, '{}'::text[]) AS allergies,
            COALESCE(play_styles, '{}'::text[]) AS play_styles
        FROM children
        WHERE user_id = $1
        LIMIT 1`, id)
	
    err := row.Scan(
        &c.UserID,
        &c.Name,
		&c.Birthday,
        &c.Gender,
        &c.About_short,
        &c.Interests,
        &c.Activity_level,
        &c.Limitations,
		&c.Allergies,
		&c.Play_styles,
    )
    if err != nil {
        log.Println(err)
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil //not found
        }
        return nil, err
    }

    return &c, nil
}
func GetUserMatchingPreferences(ctx context.Context, pool *pgxpool.Pool, userID string) (*structs.PreferencesInput, error) {
	var p structs.PreferencesInput
	row := pool.QueryRow(ctx, `
	SELECT
		user_id::text,
		COALESCE(interests_weight, 3),
		COALESCE(activity_level_weight, 3),
		COALESCE(limitations_weight, 2),
		COALESCE(allergies_weight, 2),
		COALESCE(play_styles_weight, 3),
		COALESCE(max_age_difference, 24)
	FROM matching_preferences
	WHERE user_id = $1
	LIMIT 1
	`, userID)
	
	err := row.Scan(
		&p.UserID,
		&p.InterestsWeight,
		&p.ActivityLevelWeight,
		&p.LimitationsWeight,
		&p.AllergiesWeight,
		&p.PlayStylesWeight,
		&p.MaxAgeDifference,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Create default preferences if not found
			log.Printf("No preferences found for user %s, creating defaults", userID)
			return createDefaultMatchingPreferences(ctx, pool, userID)
		}
		log.Printf("Error getting preferences: %v", err)
		return nil, err
	}
	
	return &p, nil
}

// createDefaultMatchingPreferences creates default matching preferences for a user
func createDefaultMatchingPreferences(ctx context.Context, pool *pgxpool.Pool, userID string) (*structs.PreferencesInput, error) {
	var p structs.PreferencesInput
	
	err := pool.QueryRow(ctx, `
		INSERT INTO matching_preferences 
		(user_id, interests_weight, activity_level_weight, limitations_weight, 
		 allergies_weight, play_styles_weight, max_age_difference)
		VALUES ($1, 3, 3, 2, 2, 3, 24)
		ON CONFLICT (user_id) DO UPDATE SET
			interests_weight = EXCLUDED.interests_weight
		RETURNING user_id::text, interests_weight, activity_level_weight, 
		          limitations_weight, allergies_weight, play_styles_weight, 
		          max_age_difference
	`, userID).Scan(
		&p.UserID,
		&p.InterestsWeight,
		&p.ActivityLevelWeight,
		&p.LimitationsWeight,
		&p.AllergiesWeight,
		&p.PlayStylesWeight,
		&p.MaxAgeDifference,
	)
	
	if err != nil {
		log.Printf("ERROR creating default preferences: %v", err)
		return nil, err
	}
	
	log.Printf("Created default preferences for user %s", userID)
	return &p, nil
}