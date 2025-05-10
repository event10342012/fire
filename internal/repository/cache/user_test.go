package cache

import (
	"context"
	"fire/internal/domain"
	"testing"
)

func TestRedisUserCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		user    domain.User
		wantErr error
	}{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
