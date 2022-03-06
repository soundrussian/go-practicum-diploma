package api

import (
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	validRequest = `
		{
			"login": "john.doe@example.com",
			"password": "topsecret"
		}
	`
	invalidJSON = `not a json`
)

func TestHandleRegister(t *testing.T) {
	type args struct {
		body    string
		headers map[string]string
		auth    auth.Auth
	}
	type want struct {
		status  int
		headers map[string]string
		body    string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "returns 415 Unsupported Media Type if request does not have Content-Type: application/json",
			args: args{
				body: validRequest,
				auth: &mocks.Auth{},
			},
			want: want{
				status: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "returns 400 Bad Request if request payload is not JSON",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: invalidJSON,
				auth: &mocks.Auth{},
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "returns 400 Bad Request if Auth service could not register user",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: validRequest,
				auth: invalidLogin(),
			},
			want: want{
				status: http.StatusBadRequest,
				body:   auth.ErrInvalidLogin.Error() + "\n",
			},
		},
		{
			name: "returns 409 Bad Request if user is already registered",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: validRequest,
				auth: duplicateUser(),
			},
			want: want{
				status: http.StatusConflict,
				body:   auth.ErrUserAlreadyRegistered.Error() + "\n",
			},
		},
		{
			name: "returns 200 OK if user is registered, along with authentication token",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: validRequest,
				auth: successfulRegistration(100),
			},
			want: want{
				status: http.StatusOK,
				body:   fmt.Sprintf(`{"token":"%s"}`+"\n", token(100)),
				headers: map[string]string{
					"Set-Cookie":     fmt.Sprintf("jwt=%s", token(100)),
					"Authentication": fmt.Sprintf("Bearer %s", token(100)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(tt.args.auth, new(mocks.Balance), new(mocks.Order))
			require.NoError(t, err)

			r := a.routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(tt.args.body))
			require.NoError(t, err)

			for header, value := range tt.args.headers {
				req.Header.Set(header, value)
			}

			transport := http.Transport{}
			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)

			if tt.want.body != "" {
				resBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.want.body, string(resBody))
			}

			assert.Equal(t, tt.want.status, resp.StatusCode)

			for wantHeader, wantHeaderValue := range tt.want.headers {
				assert.Equal(t, wantHeaderValue, resp.Header.Get(wantHeader))
			}

		})
	}
}

func invalidLogin() *mocks.Auth {
	m := new(mocks.Auth)
	m.On("Register", mock.Anything, mock.Anything, mock.Anything).Return(nil, auth.ErrInvalidLogin)
	return m
}

func duplicateUser() *mocks.Auth {
	m := new(mocks.Auth)
	m.On("Register", mock.Anything, mock.Anything, mock.Anything).Return(nil, auth.ErrUserAlreadyRegistered)
	return m
}

func successfulRegistration(userID uint64) *mocks.Auth {
	m := new(mocks.Auth)
	m.On("Register", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{ID: userID}, nil)
	t := token(100)
	m.On("AuthToken", mock.Anything, mock.Anything, mock.Anything).Return(&t, nil)
	return m
}
