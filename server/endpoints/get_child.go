package endpoints

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/helpers"
	"matchme-server/internal"
	"matchme-server/structs"
	"time"

	"github.com/gin-gonic/gin"
)

type ChildProfileResponse struct {
	Name                string    `json:"name"`
	Birthday            time.Time `json:"birthday"`
	Gender              string    `json:"gender"`
	About_short         string    `json:"about_short"`
	Interests           []string  `json:"interests"`
	Activity_level      string    `json:"activity_level"`
	Limitations         []string  `json:"limitations"`
	Allergies           []string  `json:"allergies"`
	Play_styles         []string  `json:"play_styles"`
	InterestsWeight     int       `json:"interests_weight,omitempty"`
	ActivityLevelWeight int       `json:"activity_level_weight,omitempty"`
	LimitationsWeight   int       `json:"limitations_weight,omitempty"`
	AllergiesWeight     int       `json:"allergies_weight,omitempty"`
	PlayStylesWeight    int       `json:"play_styles_weight,omitempty"`
	MaxAgeDifference    int       `json:"max_age_difference,omitempty"`
}

func GetChildProfile(c *gin.Context) {
	id := c.GetString("userID")
	if !helpers.IsValidID(id) {
		log.Println(id)
		c.JSON(400, structs.ErrorResponse{
			Message: "invalid id",
		})
		return
	}

	ctx := context.Background()
	ch, err := database.GetChildProfile(ctx, internal.DB, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return
	}

	pref, err := database.GetUserMatchingPreferences(ctx, internal.DB, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return
	}

	if ch == nil {
		ch = &structs.Child{}
	}
	if pref == nil {
		pref = &structs.PreferencesInput{}
	}
	res := buildChildResponse(ch, pref)
	c.JSON(200, res)
}

func buildChildResponse(ch *structs.Child, pref *structs.PreferencesInput) ChildProfileResponse {

	return ChildProfileResponse{

		Name:                ch.Name,
		Birthday:            ch.Birthday,
		Gender:              ch.Gender,
		About_short:         ch.About_short,
		Interests:           ch.Interests,
		Activity_level:      ch.Activity_level,
		Limitations:         ch.Limitations,
		Allergies:           ch.Allergies,
		Play_styles:         ch.Play_styles,
		InterestsWeight:     pref.InterestsWeight,
		ActivityLevelWeight: pref.ActivityLevelWeight,
		LimitationsWeight:   pref.LimitationsWeight,
		AllergiesWeight:     pref.AllergiesWeight,
		PlayStylesWeight:    pref.PlayStylesWeight,
		MaxAgeDifference:    pref.MaxAgeDifference,
	}
}
