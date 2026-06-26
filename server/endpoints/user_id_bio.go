package endpoints

import (
	"matchme-server/structs"
	"time"

	"github.com/gin-gonic/gin"
)

type ChildBio struct {
	Birthday       time.Time `json:"birthday"`
	Gender         string    `json:"gender"`
	Interests      []string  `json:"intersts"`
	Activity_level string    `json:"activity_level"`
	Limitations    []string  `json:"limitations"`
	Allergies      []string  `json:"allergies"`
	Play_styles    []string  `json:"play_styles"`
}

type BioRespond struct {
	ID                string   `json:"id"`
	Gender            string   `json:"gender"`
	Languages         []string `json:"languages"`
	AddressCity       string   `json:"addressCity"`
	Lat               float64  `json:"lat"`
	Lon               float64  `json:"lon"`
	PreferredDistance int      `json:"preferredDistance"`
	Child             ChildBio `json:"child"`
}

// GET /users/:id (bio)
func GetUserBioByID(c *gin.Context) {
	serveBio(c, c.Param("id"))
}

// GET /me (bio) — protect this route with AuthRequired()
func GetMeBio(c *gin.Context) {
	id := c.GetString("userID") // guaranteed by AuthRequired()
	serveBio(c, id)
}

func serveBio(c *gin.Context, id string) {
	p, ch, ok := LoadProfiles(c, id)

	if !ok {
		return 
	}
	res := buildBioResponse(p, ch)
	c.JSON(200, res)
}

func buildBioResponse(p *structs.ParentProfile, ch *structs.Child) BioRespond {
	return BioRespond{
		ID:                p.UserID,
		Gender:            p.Gender,
		Languages:         p.Languages,
		AddressCity:       p.AddressCity,
		Lat:               p.Lat,
		Lon:               p.Lon,
		PreferredDistance: p.PreferredDistance,
		Child: ChildBio{
			Birthday:       ch.Birthday,
			Gender:         ch.Gender,
			Interests:      ch.Interests,
			Activity_level: ch.Activity_level,
			Limitations:    ch.Limitations,
			Allergies:      ch.Allergies,
			Play_styles:    ch.Play_styles,
		},
	}
}
