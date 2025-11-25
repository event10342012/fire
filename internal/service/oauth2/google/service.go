package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fire/internal/domain"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	userInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type AuthService struct {
	Config *oauth2.Config
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

func (a *AuthService) FetchUserInfo(ctx context.Context, client *http.Client) (domain.GoogleUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoUrl, nil)
	if err != nil {
		return domain.GoogleUser{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return domain.GoogleUser{}, err
	}
	defer resp.Body.Close()
	var user domain.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return domain.GoogleUser{}, err
	}
	return user, nil
}
