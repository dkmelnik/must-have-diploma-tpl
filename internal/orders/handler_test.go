package orders

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders/mocks"
)

type (
	testPrepareFunc func(*servicesMock)
	servicesMock    struct {
		orderService *mocks.MockorderService
	}
	want struct {
		code        int
		contentType string
	}
	testCase struct {
		name     string
		prepare  testPrepareFunc
		body     string
		method   string
		want     want
		wantErr  bool
		response []dto.OrderResponse
	}
)

func Test_create(t *testing.T) {
	tests := []testCase{
		{
			name:    "negative test #1, bad entity: empty body",
			body:    "",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #2, bad entity: incorrect body",
			body:    "test",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #3, bad entity: check number on luhn",
			body:    "123456789",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, unknown service error",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().CreateIfNotExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(apperrors.ErrIsExist).AnyTimes()
			},
			body:    "555011422373148",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #5: status conflict",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().CreateIfNotExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("order already exists")).AnyTimes()
			},
			body:    "555011422373148",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #6: status ok",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().CreateIfNotExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("order already exists for the user")).AnyTimes()
			},
			body:    "555011422373148",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},

		{
			name: "positive test #7: status accepted",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().CreateIfNotExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			body:    "555011422373148",
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusAccepted,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := mockAndRegisterHandlers(t, tt.prepare)
			defer ts.Shutdown()

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")

			resp, err := ts.Test(req, 100)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func Test_getAllOrders(t *testing.T) {
	tests := []testCase{
		{
			name: "negative test #1: unknown service error",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().GetAllUserOrders(gomock.Any(), gomock.Any()).Return(nil, errors.New("something wrong")).AnyTimes()
			},
			method:  http.MethodGet,
			wantErr: true,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #2: status no content",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().GetAllUserOrders(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			method:  http.MethodGet,
			wantErr: false,
			want: want{
				code:        http.StatusNoContent,
				contentType: "",
			},
		},
		{
			name: "positive test #3: status ok",
			prepare: func(f *servicesMock) {
				f.orderService.EXPECT().GetAllUserOrders(gomock.Any(), gomock.Any()).Return([]dto.OrderResponse{
					{Number: "12323", Status: "PROCESSED", Accrual: new(float64), UploadedAt: time.Now()},
				}, nil).AnyTimes()
			},
			method:  http.MethodGet,
			wantErr: false,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := mockAndRegisterHandlers(t, tt.prepare)
			defer ts.Shutdown()

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")

			resp, err := ts.Test(req, 100)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func mockAndRegisterHandlers(t *testing.T, prepare testPrepareFunc) *fiber.App {
	app := fiber.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var os orderService

	if prepare != nil {
		f := servicesMock{
			mocks.NewMockorderService(ctrl),
		}
		prepare(&f)
		os = f.orderService
	}

	h := newHandler(os)
	a := func(c *fiber.Ctx) error {
		c.Locals("user_id", "test_user_id")
		return c.Next()
	}
	app.Post("/", a, h.create)
	app.Get("/", a, h.getAllOrders)

	return app
}
