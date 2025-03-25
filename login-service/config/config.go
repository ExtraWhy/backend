package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
  "log"
  "github.com/joho/godotenv"
)

var GoogleOAuthConfig *oauth2.Config

func init() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  log.Println("google url" , os.Getenv("GOOGLE_REDIRECT_URL"), os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"))

  GoogleOAuthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
    Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
    Endpoint:     google.Endpoint,
  }
}

