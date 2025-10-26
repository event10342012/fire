package auth

import (
	"context"
	"fire/internal/service/sms"

	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

func (s *SMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SMSClaim
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, claims.Tpl, args, numbers...)
}

type SMSClaim struct {
	jwt.RegisteredClaims
	Tpl string
}
