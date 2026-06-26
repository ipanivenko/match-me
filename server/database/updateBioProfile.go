package database

import (
	"context"
	"errors"
	"matchme-server/graph/model"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UpdateProfileInput struct {
	// parent_profiles
	Name        *string
	About       *string
	Languages   []string
	AddressCity *string
	Lat         *float64
	Lon         *float64

	// children
	ChildName      *string
	ChildAbout     *string // maps to children.about_short
	ChildInterests []string
}

type UpdateBioInput struct {
	// parent_profiles
	ParentGender      *model.GenderEnum // maps to parent_profiles.gender
	PreferredDistance *int32    // maps to parent_profiles.preferred_distance_km

	// children
	ChildBirthday      *string // expected format: "YYYY-MM-DD"
	ChildGender        *model.ChidGenderEnum 
	ChildActivityLevel *string
	Limitations        []string
	Allergies          []string
	PlayStyles         []string
}


type setBuilder struct {
	cols []string
	args []any
}

func (b *setBuilder) add(col string, val any) {
	b.cols = append(b.cols, col)
	b.args = append(b.args, val)
}

func (b *setBuilder) empty() bool { return len(b.cols) == 0 }

// ---------- UpdateProfile: updates parent_profiles + children ----------

func UpdateProfile(ctx context.Context, db *pgxpool.Pool, userID string, in UpdateProfileInput) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // safe rollback

	// --- parent_profiles
	p := setBuilder{}
	if in.Name != nil {
		p.add("name = $"+strconv.Itoa(len(p.args)+1), *in.Name)
	}
	if in.About != nil {
		p.add("about = $"+strconv.Itoa(len(p.args)+1), *in.About)
	}
	if in.Languages != nil {
		p.add("languages = $"+strconv.Itoa(len(p.args)+1), in.Languages)
	}
	if in.AddressCity != nil {
		p.add("address_city = $"+strconv.Itoa(len(p.args)+1), *in.AddressCity)
	}
	if in.Lat != nil {
		p.add("lat = $"+strconv.Itoa(len(p.args)+1), *in.Lat)
	}
	if in.Lon != nil {
		p.add("lon = $"+strconv.Itoa(len(p.args)+1), *in.Lon)
	}

	if !p.empty() {
		p.cols = append(p.cols, "updated_at = now()")

		query := "UPDATE parent_profiles SET " + joinComma(p.cols) + " WHERE user_id = $" + strconv.Itoa(len(p.args)+1)
		p.args = append(p.args, userID)
		ct, execErr := tx.Exec(ctx, query, p.args...)
		if execErr != nil {
			return execErr
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no parent_profiles row found for user_id")
		}
	}

	// --- children
	c := setBuilder{}
	if in.ChildName != nil {
		c.add("name = $"+strconv.Itoa(len(c.args)+1), *in.ChildName)
	}
	if in.ChildAbout != nil {
		c.add("about_short = $"+strconv.Itoa(len(c.args)+1), *in.ChildAbout)
	}
	if in.ChildInterests != nil {
		c.add("interests = $"+strconv.Itoa(len(c.args)+1), in.ChildInterests)
	}

	if !c.empty() {
		query := "UPDATE children SET " + joinComma(c.cols) + " WHERE user_id = $" + strconv.Itoa(len(c.args)+1)
		c.args = append(c.args, userID)
		ct, execErr := tx.Exec(ctx, query, c.args...)
		if execErr != nil {
			return execErr
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no children row found for user_id")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}



func UpdateBio(ctx context.Context, db *pgxpool.Pool, userID string, in UpdateBioInput) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	p := setBuilder{}
	if in.ParentGender != nil {
		p.add("gender = $"+strconv.Itoa(len(p.args)+1), *in.ParentGender)
	}
	if in.PreferredDistance != nil {
		p.add("preferred_distance_km = $"+strconv.Itoa(len(p.args)+1), *in.PreferredDistance)
	}
	if !p.empty() {
		p.cols = append(p.cols, "updated_at = now()")
		query := "UPDATE parent_profiles SET " + joinComma(p.cols) + " WHERE user_id = $" + strconv.Itoa(len(p.args)+1)
		p.args = append(p.args, userID)
		ct, execErr := tx.Exec(ctx, query, p.args...)
		if execErr != nil {
			return execErr
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no parent_profiles row found for user_id")
		}
	}

	c := setBuilder{}
	if in.ChildBirthday != nil {
		t, parseErr := time.Parse("2006-01-02", *in.ChildBirthday)
		if parseErr != nil {
			return parseErr
		}
		c.add("birthday = $"+strconv.Itoa(len(c.args)+1), t)
	}
	if in.ChildGender != nil {
		c.add("gender = $"+strconv.Itoa(len(c.args)+1), *in.ChildGender)
	}
	if in.ChildActivityLevel != nil {
		c.add("activity_level = $"+strconv.Itoa(len(c.args)+1), *in.ChildActivityLevel)
	}
	if in.Limitations != nil {
		c.add("limitations = $"+strconv.Itoa(len(c.args)+1), in.Limitations)
	}
	if in.Allergies != nil {
		c.add("allergies = $"+strconv.Itoa(len(c.args)+1), in.Allergies)
	}
	if in.PlayStyles != nil {
		c.add("play_styles = $"+strconv.Itoa(len(c.args)+1), in.PlayStyles)
	}

	if !c.empty() {
		query := "UPDATE children SET " + joinComma(c.cols) + " WHERE user_id = $" + strconv.Itoa(len(c.args)+1)
		c.args = append(c.args, userID)
		ct, execErr := tx.Exec(ctx, query, c.args...)
		if execErr != nil {
			return execErr
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no children row found for user_id")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}



func joinComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += ", " + parts[i]
	}
	return out
}
