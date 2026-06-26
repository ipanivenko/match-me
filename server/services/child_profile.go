package services

import (
	"errors"
	"fmt"
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"
	"github.com/gin-gonic/gin"
)

// func for updating or filling the profile
// handlers/me_profile.go

type PatchChildProfileInput struct {
	Name           *string    `json:"name,omitempty"`
	Birthday       *string `json:"birthday,omitempty"`
	Gender         *string    `json:"gender,omitempty"`
	About_short    *string    `json:"about_short,omitempty"`
	Interests      *[]string  `json:"interests"`
	Activity_level *string    `json:"activity_level"`
	Limitations    *[]string  `json:"limitations"`
	Allergies      *[]string  `json:"allergies"`
	Play_styles    *[]string  `json:"play_styles"`
	InterestsWeight     *int `json:"interests_weight,omitempty"`
	ActivityLevelWeight *int `json:"activity_level_weight,omitempty"`
	LimitationsWeight   *int `json:"limitations_weight,omitempty"`
	AllergiesWeight     *int `json:"allergies_weight,omitempty"`
	PlayStylesWeight    *int `json:"play_styles_weight,omitempty"`
	MaxAgeDifference    *int `json:"max_age_difference,omitempty"`
}

func PatchMeChild(c *gin.Context) {
	uid := c.GetString("userID") // set by middleware

	var in PatchChildProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		log.Println(err)
		c.JSON(400, structs.ErrorResponse{Message: "invalid json"})
		return
	}

	
	childSets := []string{}
	childArgs := []any{uid}
	childUpdatedCols := make([]string, 0, 10)

	addChild := func(col string, v any) {
		childSets = append(childSets, fmt.Sprintf("%s=$%d", col, len(childArgs)+1))
		childArgs = append(childArgs, v)
		childUpdatedCols = append(childUpdatedCols, col)
	}

	if in.Name != nil {
		addChild("name", *in.Name)
	}
	if in.Birthday != nil {
		addChild("birthday", *in.Birthday)
	}
	if in.Gender != nil {
		addChild("gender", *in.Gender)
	}
	if in.About_short != nil {
		addChild("about_short", *in.About_short)
	}
	if in.Interests != nil {
		addChild("interests", *in.Interests) 
	}
	if in.Activity_level != nil {
		addChild("activity_level", *in.Activity_level)
	}
	if in.Limitations != nil {
		addChild("limitations", *in.Limitations)
	}
	if in.Allergies != nil {
		addChild("allergies", *in.Allergies)
	}
	if in.Play_styles != nil {
		addChild("play_styles", *in.Play_styles)
	}


	prefSets := []string{}
	prefArgs := []any{uid}
	prefUpdatedCols := make([]string, 0, 8)

	addPref := func(col string, v any) {
		prefSets = append(prefSets, fmt.Sprintf("%s=$%d", col, len(prefArgs)+1))
		prefArgs = append(prefArgs, v)
		prefUpdatedCols = append(prefUpdatedCols, col)
	}


	if in.InterestsWeight != nil {
		addPref("interests_weight", *in.InterestsWeight)
	}
	if in.ActivityLevelWeight != nil {
		addPref("activity_level_weight", *in.ActivityLevelWeight)
	}
	if in.LimitationsWeight != nil {
		addPref("limitations_weight", *in.LimitationsWeight)
	}
	if in.AllergiesWeight != nil {
		addPref("allergies_weight", *in.AllergiesWeight)
	}
	if in.PlayStylesWeight != nil {
		addPref("play_styles_weight", *in.PlayStylesWeight)
	}
	if in.MaxAgeDifference != nil {
		addPref("max_age_difference", *in.MaxAgeDifference)
	}

	if len(childSets) == 0 && len(prefSets) == 0 {
		c.JSON(400, structs.ErrorResponse{Message: "no fields to update"})
		return
	}

	ctx := c.Request.Context()

	var childMap map[string]any
	var prefMap map[string]any

	// update children (if any fields)
	if len(childSets) > 0 {
		m, err := database.UpdateProfileDynamic(
			ctx,
			internal.DB,
			childSets,
			childArgs,
			childUpdatedCols,
			"children",
		)
		if err != nil {
			log.Println(err)
			if errors.Is(err, database.ErrProfileNotFound) {
				c.JSON(404, structs.ErrorResponse{Message: "child profile not found"})
				return
			}
			c.JSON(500, structs.ErrorResponse{Message: CommonErr})
			return
		}
		childMap = m
	}

	// update matching_preferences (if any fields)
	if len(prefSets) > 0 {
		m, err := database.UpdateProfileDynamic(
			ctx,
			internal.DB,
			prefSets,
			prefArgs,
			prefUpdatedCols,
			"matching_preferences", 
		)
		if err != nil {
			log.Println(err)
			if errors.Is(err, database.ErrProfileNotFound) {
				// preferences row missing: surface a clear 404 to client
				c.JSON(404, structs.ErrorResponse{Message: "matching preferences not found"})
				return
			}
			c.JSON(500, structs.ErrorResponse{Message: CommonErr})
			return
		}
		prefMap = m
	}

	// response: keep backward-compat (single map) if only one part updated
	if childMap != nil && prefMap == nil {
		c.JSON(200, childMap)
		return
	}
	if prefMap != nil && childMap == nil {
		c.JSON(200, prefMap)
		return
	}
	// both updated — return grouped payload
	c.JSON(200, gin.H{
		"child":                 childMap,
		"matching_preferences":  prefMap,
	})
}
