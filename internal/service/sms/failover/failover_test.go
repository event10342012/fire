package failover

import (
	"context"
	"errors"
	"fire/internal/service/sms"
	smsmocks "fire/internal/service/sms/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "one time success",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return []sms.Service{svc}
			},
		},
		{
			name: "second time success",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("send failed"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
		},
		{
			name: "second time success",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("send failed"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("send failed"))
				return []sms.Service{svc0, svc1}
			},
			wantErr: errors.New("all SMS service failed to send"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsServices := tc.mock(ctrl)
			smsSvc := NewFailOverSMSService(smsServices)
			err := smsSvc.Send(context.Background(), "", []string{}, "")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
