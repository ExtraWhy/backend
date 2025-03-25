package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"login-service/config"
	"slices"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var allowedRedirectHosts []string

func init() {
  hosts := os.Getenv("ALLOWED_HOSTS")
  if hosts != "" {
    allowedRedirectHosts = strings.Split(hosts, ",")
  }
}

func isSafeRedirect(rawURL string) bool {
  parsedURL, err := url.Parse(rawURL)

  if ( err != nil || parsedURL.Host == "" ){
    return false
  }

  origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
  return slices.Contains(allowedRedirectHosts, origin)
}

func GoogleLogin(c *gin.Context) {
  state := c.Query("callback-url")
  url := config.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	client := config.GoogleOAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo map[string]any
	json.Unmarshal(body, &userInfo)

	// Now from here we can Look up the user in DB, issue JWT and/or create a session

  callback_url := c.Query("state") 
  if ( !isSafeRedirect(callback_url) ) {
    fmt.Println("NOT A SAFE REDIRECT DICKHEAD origin", callback_url)
    c.AbortWithStatus(http.StatusBadRequest)
    // c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid redirect"})
    return
  }

  familyName, _ := userInfo["family_name"].(string)
  givenName, _ := userInfo["given_name"].(string)
  picture, _ := userInfo["picture"].(string)

  c.SetCookie("token", token.AccessToken, 3600, "/", callback_url, true, true)
  c.SetCookie("photo", picture, 3600, "/", callback_url, true, false)
  c.SetCookie("username", familyName + " " + givenName, 3600, "/", callback_url, true, false)

  c.Redirect(http.StatusTemporaryRedirect, callback_url)
}
