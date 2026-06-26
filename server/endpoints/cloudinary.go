package endpoints

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
	"matchme-server/internal"

	"github.com/gin-gonic/gin"
)

// GET /api/me/cloudinary-sign
func CloudinarySign(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatus(401)
		return
	}

	folder := "users/" + userID
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	publicID := "avatar"  
	overwrite := "true"   

	raw := "folder=" + folder +
		"&overwrite=" + overwrite +
		"&public_id=" + publicID +
		"&timestamp=" + ts +
		internal.Cfg.Cloud_secret

	sum := sha1.Sum([]byte(raw))
	sig := hex.EncodeToString(sum[:])

	c.JSON(200, gin.H{
		"cloud_name": internal.Cfg.Cloud_name,
		"api_key":    internal.Cfg.Cloud_key,
		"timestamp":  ts,
		"signature":  sig,
		"folder":     folder,
		"public_id":  publicID,   
		"overwrite":  overwrite, 
	})
}
