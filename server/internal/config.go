// file for loading server/.env and saving it into struct in order to use after
package internal

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	IsDevMode   bool//for graphql
	Port        string
	DatabaseURL string
	JWTSecret   string
	Cloud_secret string
	Cloud_name string
	Cloud_key string

}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	} else {
		log.Println(".env read successfully")
	}

	cloud := strings.TrimSpace(os.Getenv("CLOUDINARY_CLOUD_NAME"))
	if cloud == "" {
		log.Println("CLOUDINARY_CLOUD_NAME is empty")
	}

	var c Config
	mode := os.Getenv("MODE")
	if mode == "development"{
		c.IsDevMode = true
	}
	c.Port = os.Getenv("PORT")
	c.DatabaseURL = os.Getenv("DATABASE_URL")
	c.JWTSecret = os.Getenv("JWT_SECRET")
	c.Cloud_secret = os.Getenv("CLOUDINARY_API_SECRET")
	c.Cloud_name = cloud
	c.Cloud_key = os.Getenv("CLOUDINARY_API_KEY")
	Cfg = &c
	return Cfg
}
