package users

import (
	"encoding/json"
	"errors"
	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dkmelnik/go-musthave-diploma/internal/users/mocks"
)

type (
	testPrepareFunc func(*servicesMock)
	servicesMock    struct {
		userService *mocks.MockuserService
	}
	want struct {
		code        int
		contentType string
	}
	testCase struct {
		name    string
		prepare testPrepareFunc
		body    map[string]interface{}
		method  string
		want    want
		wantErr bool
	}
)

func Test_registration(t *testing.T) {
	tests := []testCase{
		{
			name:    "negative test #1, bad entity: empty body",
			body:    nil,
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2, bad entity: invalid login field type",
			body: map[string]interface{}{
				"login":    12213123123,
				"password": "testtest21@",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3, bad entity: invalid password field type",
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": 12213123123,
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, login exist",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", apperrors.ErrIsExist).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #5, unknown service error",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", errors.New("some error")).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #6, user saved, token returned",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("token", nil).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: false,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := mockAndRegisterHandlers(t, tt.prepare)
			defer ts.Shutdown()

			bts, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(tt.method, "/register", strings.NewReader(string(bts)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Test(req, 100)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			if !tt.wantErr {
				dump := utils.DumpRequest(req, true)
				setCookieHeader := resp.Header.Get("Set-Cookie")
				assert.True(t, setCookieHeader != "",
					"Не удалось обнаружить авторизационные данные в ответе", string(dump))
			}
		})
	}
}

func Test_authenticate(t *testing.T) {
	tests := []testCase{
		{
			name:    "negative test #1, bad entity: empty body",
			body:    nil,
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2, bad entity: invalid login field type",
			body: map[string]interface{}{
				"login":    12213123123,
				"password": "testtest21@",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3, bad entity: invalid password field type",
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": 12213123123,
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, login not exist",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return("", apperrors.ErrNotFound).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #5, invalid password",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return("", apperrors.ErrInvalidCredentials).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #6, unknown service error",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return("", errors.New("some error")).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #7, authenticate success, token returned",
			prepare: func(f *servicesMock) {
				f.userService.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return("token", nil).AnyTimes()
			},
			body: map[string]interface{}{
				"login":    "testtest21@",
				"password": "12213123123",
			},
			method:  http.MethodPost,
			wantErr: false,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := mockAndRegisterHandlers(t, tt.prepare)
			defer ts.Shutdown()

			bts, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(string(bts)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Test(req, 100)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			if !tt.wantErr {
				dump := utils.DumpRequest(req, true)
				setCookieHeader := resp.Header.Get("Set-Cookie")
				assert.True(t, setCookieHeader != "",
					"Не удалось обнаружить авторизационные данные в ответе", string(dump))
			}
		})
	}
}

func mockAndRegisterHandlers(t *testing.T, prepare testPrepareFunc) *fiber.App {
	app := fiber.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var us userService

	if prepare != nil {
		f := servicesMock{
			mocks.NewMockuserService(ctrl),
		}
		prepare(&f)
		us = f.userService
	}

	h := newHandler(time.Hour, us)
	app.Post("/register", h.register)
	app.Post("/login", h.authenticate)

	return app
}
