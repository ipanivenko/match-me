package graph

import (
	"context"
	"errors"
	"log"
	"matchme-server/database"
	"matchme-server/graph/model"
	"matchme-server/structs"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *pgxpool.Pool
}

var GlobalPubSub = NewPubSub()

func stringSliceToPtrSlice(s []string) []*string {
	if s == nil {
		return nil
	}
	ptrSlice := make([]*string, len(s))
	for i := range s {
		ptrSlice[i] = &s[i]
	}
	return ptrSlice
}

func (r *Resolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	email, created_at, err := database.GetUserEmailCreatedAt(ctx, r.DB, id)
	if err != nil {
		log.Printf("Error in GetUserEmailCreatedAt for ID %s: %v", id, err) // Log the actual error
		return nil, gqlerror.Errorf("Database error fetching user details")
	}

	_, profilePicture := database.GetUserNamePhotoURL(ctx, r.DB, id)

	user := &model.User{
		UserID:         id,
		Email:          email,
		CreatedAt:      created_at,
		ProfilePicture: &profilePicture,
	}

	return user, nil
}

func (r *Resolver) GetProfile(ctx context.Context, userID string) (*model.Profile, error) {
	p, ch, err := r.LoadProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var latPtr *float64
	if p.Lat != 0.0 {
		temp := float64(p.Lat)
		latPtr = &temp
	}

	var lonPtr *float64
	if p.Lon != 0.0 {
		temp := float64(p.Lon)
		lonPtr = &temp
	}

	// Create the final response object
	profile := &model.Profile{
		UserID:         p.UserID,
		Name:           &p.Name,
		About:          &p.About,
		Languages:      stringSliceToPtrSlice(p.Languages),
		AddressCity:    &p.AddressCity,
		Lat:            latPtr,
		Lon:            lonPtr,
		ChildName:      &ch.Name,
		ChildAbout:     &ch.About_short,
		ChildInterests: stringSliceToPtrSlice(ch.Interests),
	}

	return profile, nil
}

func (r *Resolver) GetBio(ctx context.Context, userID string) (*model.Bio, error) {
	p, ch, err := r.LoadProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var preferredDistPtr *int32
	if p.PreferredDistance != 0 {
		temp := int32(p.PreferredDistance) // 1. Convert int to int32
		preferredDistPtr = &temp           // 2. Get a pointer to the int32
	}

	var birthdayStrPtr *string

	if !ch.Birthday.IsZero() {
		tempStr := ch.Birthday.Format("2006-01-02")
		birthdayStrPtr = &tempStr
	}

	bio := &model.Bio{
		UserID:             p.UserID,
		ParentGender:       model.GenderEnum(p.Gender),
		PreferredDistance:  preferredDistPtr,
		ChildBirthday:      birthdayStrPtr,
		ChildGender:        model.ChidGenderEnum(ch.Gender),
		ChildActivityLevel: model.ChildActivityLevelEnum(ch.Activity_level),
		Limitations:        stringSliceToPtrSlice(ch.Limitations),
		Allergies:          stringSliceToPtrSlice(ch.Allergies),
		PlayStyles:         stringSliceToPtrSlice(ch.Play_styles),
	}

	return bio, nil
}

func (r *Resolver) LoadProfiles(ctx context.Context, userID string) (*structs.ParentProfile, *structs.Child, error) {
	// Fetch Parent data
	p, err := database.GetUserProfile(ctx, r.DB, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return a specific GQL error
			return nil, nil, gqlerror.Errorf("Profile not found for user ID: %s", userID)
		}
		// Return a generic server error
		log.Println(err)
		return nil, nil, gqlerror.Errorf("Database error fetching parent.")
	}

	// Fetch Child data
	ch, err := database.GetChildProfile(ctx, r.DB, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, gqlerror.Errorf("Child profile not found for user ID: %s", userID)
		}
		log.Println(err)
		return nil, nil, gqlerror.Errorf("Database error fetching child")
	}

	// Success: return the two models
	return p, ch, nil
}
