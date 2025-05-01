package web

import (
	"bytes"
	"fire/internal/domain"
	"fire/internal/service"
	svcmock "fire/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)
import "go.uber.org/mock/gomock"

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.CodeService, service.UserService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "Signup success",
			mock: func(ctrl *gomock.Controller) (service.CodeService, service.UserService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:       "test@gmail.com",
					Password:    "hello#world123",
					IsSuperUser: false,
					IsActive:    true,
				}).Return(nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return codeSvc, userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
{"email": "test@gmail.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"}
`)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "Signup success",
		},

		{
			name: "Signup bind error",
			mock: func(ctrl *gomock.Controller) (service.CodeService, service.UserService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return codeSvc, userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
{"email": "test@gmail.com",
"password": "hello#world123",
`)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusBadRequest,
		},

		{
			name: "Signup email invalid",
			mock: func(ctrl *gomock.Controller) (service.CodeService, service.UserService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return codeSvc, userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
{"email": "test@gmail",
"password": "hello#world123",
"confirmPassword": "hello#world123"}
`)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "email is invalid",
		},

		{
			name: "Signup password is not match",
			mock: func(ctrl *gomock.Controller) (service.CodeService, service.UserService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return codeSvc, userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
{"email": "test@gmail.com",
"password": "hello#world123",
"confirmPassword": "hello#"}
`)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "password is not match",
		},
		{
			name: "Signup password is invalid",
			mock: func(ctrl *gomock.Controller) (service.CodeService, service.UserService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return codeSvc, userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
{"email": "test@gmail.com",
"password": "hello",
"confirmPassword": "hello"}
`)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "password is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			codeSvc, userSvc := tc.mock(ctrl)
			handle := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			handle.RegisterRoutes(server)
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, recorder.Code, tc.wantCode)
			assert.Equal(t, recorder.Body.String(), tc.wantBody)
		})
	}

}
