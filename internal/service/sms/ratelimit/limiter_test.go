package ratelimit

import (
	"fire/internal/service/sms"
	smsmocks "fire/internal/service/sms/mocks"
	"fire/pkg/limiter"
	limitermocks "fire/pkg/limiter/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimiterSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		wantErr error
	}{
		{
			name: "valid sms send",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return smsSvc, l
			},
			wantErr: nil,
		},
		{
			name: "limit send error",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return smsSvc, l
			},
			wantErr: errLimited,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewRateLimiterSMSService(smsSvc, l)
			err := svc.Send(nil, "", nil, "")
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
