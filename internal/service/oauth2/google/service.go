package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

// Keys used in session map
const (
	SessionKeyState      = "oauth_state"
	SessionKeyIDToken    = "id_token"
	SessionKeyAccessTok  = "access_token"
	SessionKeyRefreshTok = "refresh_token"
	SessionKeyUser       = "user_info" // JSON string
	UserInfoUrl          = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type AuthService struct {
	Config *oauth2.Config
}

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func NewService(config *oauth2.Config) *AuthService {
	return &AuthService{
		Config: config,
	}
}

// NewState returns a URL-safe random string for OAuth2 state param.
func (a *AuthService) NewState(n int) (string, error) {
	if n <= 0 {
		n = 32
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// AuthCodeURL returns a Google OAuth2 auth URL given state.
func (a *AuthService) AuthCodeURL(state string) string {
	return a.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCode exchanges the authorization code for tokens.
func (a *AuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.Config.Exchange(ctx, code)
}

// Client returns an HTTP client authorized with the given token.
func (a *AuthService) Client(ctx context.Context, tok *oauth2.Token) *http.Client {
	return a.Config.Client(ctx, tok)
}

func (a *AuthService) FetchUserInfo(ctx context.Context, client *http.Client) (User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", UserInfoUrl, nil)
	if err != nil {
		return User{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
