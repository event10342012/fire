package server

import (
	"context"
	"encoding/json"
	"fire/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback", // This should match the one you set in the Google Developer Console
		ClientID:     "223534559756-idf2u1r7uglqokpsedlevld2djk7qmeh.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func LoginHandler(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	userinfo, err := getUserInfo(token)
	err = saveUser(&userinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save user "})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get User info"})
	}

	c.JSON(http.StatusOK, userinfo)
}

func getUserInfo(token *oauth2.Token) (UserInfo, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return UserInfo{}, err
	}

	return userInfo, nil
}

// saveUser saves the user if not exist in the database else skip
func saveUser(userInfo *UserInfo) error {
	db := GetDB()
	user := model.User{}
	result := db.Take(&user, "email = ?", userInfo.Email)
	if result.Error != nil {
		user = model.User{
			Email:       userInfo.Email,
			GivenName:   userInfo.GivenName,
			FamilyName:  userInfo.FamilyName,
			Picture:     userInfo.Picture,
			Locale:      userInfo.Locale,
			GoogleId:    userInfo.Id,
			IsSuperUser: false,
			IsActive:    true,
		}
		db.Create(&user)
	}
	return nil
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}
