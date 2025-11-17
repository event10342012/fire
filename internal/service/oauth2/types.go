package oauth2

import "context"

type AuthService interface {
	AuthURL(ctx context.Context) string
}
