package endpoints

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/helpers"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)


func LoadProfiles(c *gin.Context, id string) (*structs.ParentProfile, *structs.Child, bool) {

	if !helpers.IsValidID(id) { 
		log.Println(id)
		c.JSON(400, structs.ErrorResponse{
			Message: "invalid id",
		})
		return nil, nil, false
	}

	ctx := context.Background()
	p, err := database.GetUserProfile(ctx, internal.DB, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return nil, nil, false
	}

	if p == nil {
		p = &structs.ParentProfile{}
	}

	ch, err := database.GetChildProfile(ctx, internal.DB, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return nil, nil, false
	}

	if ch == nil {
		ch = &structs.Child{}
	}
	return p, ch, true
}
