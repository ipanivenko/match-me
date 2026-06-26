package endpoints

import (
	"matchme-server/helpers"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

type ChildRespond struct {
	Name         string   `json:"name"`
	AgeYears     int      `json:"ageYears"`
	Gender       string   `json:"gender"`
	AboutShort   string   `json:"aboutShort"`
	TopInterests []string `json:"topInterests"`
}

type ProfileRespond struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	About       string       `json:"about"`
	Languages   []string     `json:"languages"`
	AddressCity string       `json:"addressCity"`
	Child       ChildRespond `json:"child"`
}


// GET /users/:id/profile
func GetUserProfileByID(c *gin.Context) {
    serveProfile(c, c.Param("id"))
}

// GET /me/profile
func GetMyProfile(c *gin.Context) {
    id := c.GetString("userID")
    serveProfile(c, id)
}


func serveProfile(c *gin.Context, id string) {
	p, ch, ok := LoadProfiles(c, id)

	if !ok {
		return //json respond already sent in LoadProfiles func
	}

	res := buildProfileResponse(p, ch)
	c.JSON(200, res)
}


func buildProfileResponse(p *structs.ParentProfile, ch *structs.Child) ProfileRespond {

	childAge := helpers.ComputeAge(ch.Birthday)

	return ProfileRespond{
		ID:          p.UserID,
		Name:        p.Name,
		About:       p.About,
		Languages:   p.Languages,
		AddressCity: p.AddressCity,
		Child: ChildRespond{
			Name:         ch.Name,
			AgeYears:     childAge,
			Gender:       ch.Gender,
			AboutShort:   ch.About_short,
			TopInterests: ch.Interests,
		},
	}
}
