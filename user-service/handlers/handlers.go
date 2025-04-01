package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/models/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	Config              *config.UserService
	GoogleOAuthConfig   *oauth2.Config
	FacebookOAuthConfig *oauth2.Config
	dbc                 *db.DBConnection
	cookieExpiry        int
}

func (handler *OAuthHandler) Init(dbc *db.DBConnection) error {
	handler.dbc = dbc
	handler.GoogleOAuthConfig = buildOAuthConfig(handler.Config.GoogleProvider, google.Endpoint)
	handler.FacebookOAuthConfig = buildOAuthConfig(handler.Config.FacebookProvider, facebook.Endpoint)

	handler.cookieExpiry = 3600

	gin_engine := gin.Default()
	gin_engine.SetTrustedProxies(nil)

	gin_engine.GET("/auth/google/login", handler.GoogleLogin)
	gin_engine.GET("/auth/google/callback", handler.GoogleCallback)

	gin_engine.GET("/auth/facebook/login", handler.FacebookLogin)
	gin_engine.GET("/auth/facebook/callback", handler.FacebookCallback)

	gin_engine.GET("/users", handler.getUsers)

	addr := fmt.Sprintf("%s:%s", handler.Config.UserServiceHost, handler.Config.UserServicePort)
	gin_engine.Run(addr)

	return nil
}

func (handler *OAuthHandler) getUsers(ctx *gin.Context) {
	p, err := handler.dbc.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Users DB Error"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, p)
}

func buildOAuthConfig(provider config.OAuthProviderConfig, endpoint oauth2.Endpoint) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectUrl,
		Scopes:       provider.Scopes,
		Endpoint:     endpoint,
	}
}

func (handler *OAuthHandler) GoogleLogin(c *gin.Context) {
	handler.performLogin(c, handler.GoogleOAuthConfig)
}

func (handler *OAuthHandler) GoogleCallback(c *gin.Context) {
	handler.handleOAuthCallback(c, handler.GoogleOAuthConfig, handler.Config.GoogleProvider.UserInfoUrl, func(userInfo map[string]any) (string, string, string) {
		familyName, _ := userInfo["family_name"].(string)
		givenName, _ := userInfo["given_name"].(string)
		picture, _ := userInfo["picture"].(string)
		email, _ := userInfo["email"].(string)
		return fmt.Sprintf("%s %s", familyName, givenName), email, picture
	})
}

func (handler *OAuthHandler) FacebookLogin(c *gin.Context) {
	handler.performLogin(c, handler.FacebookOAuthConfig)
}

func (handler *OAuthHandler) FacebookCallback(c *gin.Context) {
	handler.handleOAuthCallback(c, handler.FacebookOAuthConfig, handler.Config.FacebookProvider.UserInfoUrl, func(userInfo map[string]any) (string, string, string) {
		firstName, _ := userInfo["first_name"].(string)
		lastName, _ := userInfo["last_name"].(string)
		email, _ := userInfo["email"].(string)
		picture := ""
		// Facebook returns the picture as a nested object.
		if picObj, ok := userInfo["picture"].(map[string]any); ok {
			if data, ok := picObj["data"].(map[string]any); ok {
				picture, _ = data["url"].(string)
			}
		}
		return fmt.Sprintf("%s %s", firstName, lastName), email, picture
	})
}

func (handler *OAuthHandler) performLogin(c *gin.Context, oauthConfig *oauth2.Config) {
	state := c.Query("callback-url")

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (handler *OAuthHandler) isSafeRedirect(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	fmt.Println(parsedURL, parsedURL.Host)
	if err != nil || parsedURL.Host == "" {
		return false
	}

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	fmt.Println(origin)
	return slices.Contains(handler.Config.AllowedHosts, origin)
}

func (handler *OAuthHandler) handleOAuthCallback(c *gin.Context, oauthConfig *oauth2.Config, userInfoURL string, extractUserData func(map[string]any) (username, photo, email string)) {
	code := c.Query("code")
	fmt.Println("trying to insert user")
	fmt.Println(code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	client := oauthConfig.Client(c, token)
	resp, err := client.Get(userInfoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo map[string]any
	json.Unmarshal(body, &userInfo)

	callbackURL := c.Query("state")
	if !handler.isSafeRedirect(callbackURL) {
		fmt.Println("NOT A SAFE REDIRECT origin", callbackURL)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	username, email, picture := extractUserData(userInfo)

	fmt.Println("trying to insert user")
	user := user.User{Username: username, Token: token.AccessToken, Email: email, Picture: picture}
	insert_err := handler.dbc.InsertUser(user)
	if insert_err != nil {
		fmt.Println("insert user err", insert_err)
	}
	fmt.Println("user should be inserted")

	c.SetCookie("token", token.AccessToken, handler.cookieExpiry, "/", callbackURL, true, true)
	c.SetCookie("photo", picture, handler.cookieExpiry, "/", callbackURL, true, false)
	c.SetCookie("username", username, handler.cookieExpiry, "/", callbackURL, true, false)

	c.Redirect(http.StatusTemporaryRedirect, callbackURL)
}
