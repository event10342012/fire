package ratelimit

import (
	"context"
	"errors"
	"fire/internal/service/sms"
	"fire/pkg/limiter"
)

var errLimited = errors.New("rate limit exceeded")

type RateLimiterSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimiterSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimiterSMSService(svc sms.Service, limiter limiter.Limiter) *RateLimiterSMSService {
	return &RateLimiterSMSService{
		svc:     svc,
		limiter: limiter,
		key:     "sms-limiter",
	}
}
