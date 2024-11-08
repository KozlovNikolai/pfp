package httpserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/mocks"
)

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := mocks.NewIUserService(t)
	h := HTTPServer{
		userService: mockUserService,
	}
	testCases := []struct {
		name               string
		user               domain.User
		queryInURL         string
		wantErr            bool
		mockGetUserByID    func()
		mockGetUserByLogin func()
		expectStatusCode   int
		expectResponseBody string
	}{
		{
			name: "OK get by ID",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?id=3",
			wantErr:    false,
			mockGetUserByID: func() {
				mockUserService.On("GetUserByID", mock.AnythingOfType(ginContext), mock.AnythingOfType("int")).
					Return(domain.NewUser(domain.NewUserData{
						ID:        3,
						Login:     loginAdmin,
						Password:  "hashpassword",
						Role:      "admin",
						Token:     "",
						CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
						UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
					}), nil).
					Once()
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusOK,
			expectResponseBody: `{"id":3,"login":"cmd@cmd.ru","password":"hashpassword","role":"admin","token":""}`,
		},
		{
			name: "OK get by Login",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?login=cmd@cmd.ru",
			wantErr:    false,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() {
				mockUserService.On("GetUserByLogin", mock.AnythingOfType(ginContext), mock.AnythingOfType("string")).
					Return(domain.NewUser(domain.NewUserData{
						ID:        3,
						Login:     loginAdmin,
						Password:  "hashpassword",
						Role:      "admin",
						Token:     "",
						CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
						UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
					}), nil).
					Once()
			},
			expectStatusCode:   http.StatusOK,
			expectResponseBody: `{"id":3,"login":"cmd@cmd.ru","password":"hashpassword","role":"admin","token":""}`,
		},
		{
			name: "not user in context",
			user: domain.User{ // empty value
			},
			queryInURL: "?login=cmd@cmd.ru",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `{"error":"no user in context"}`,
		},
		{
			name: "invalid login",
			user: domain.User{ // empty value
			},
			queryInURL: "?login=cmd@cmdru",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `validation`,
		},
		{
			name: "invalid zero id",
			user: domain.User{ // empty value
			},
			queryInURL: "?id=0",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `id lower or equal zero`,
		},
		{
			name: "invalid type id",
			user: domain.User{ // empty value
			},
			queryInURL: "?id=any",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `invalid-user-id`,
		},
		{
			name: "not auth get by Login",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "regular",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?login=cmd@cmd.com",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() {
				mockUserService.On("GetUserByLogin", mock.AnythingOfType(ginContext), mock.AnythingOfType("string")).
					Return(domain.NewUser(domain.NewUserData{
						ID:        4,
						Login:     loginRegular,
						Password:  "hashpassword",
						Role:      "regular",
						Token:     "",
						CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
						UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
					}), nil).
					Once()
			},
			expectStatusCode:   http.StatusUnauthorized,
			expectResponseBody: `{"invalid user login or role":"access denied"}`,
		},
		{
			name: "not auth get by ID",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "regular",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?id=4",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusUnauthorized,
			expectResponseBody: `{"invalid user id or role":"access denied"}`,
		},
		{
			name: "admin get by ID",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?id=4",
			wantErr:    false,
			mockGetUserByID: func() {
				mockUserService.On("GetUserByID", mock.AnythingOfType(ginContext), mock.AnythingOfType("int")).
					Return(domain.NewUser(domain.NewUserData{
						ID:        4,
						Login:     loginRegular,
						Password:  "hashpassword",
						Role:      "admin",
						Token:     "",
						CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
						UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
					}), nil).
					Once()
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusOK,
			expectResponseBody: `{"id":4,"login":"cmd@cmd.com","password":"hashpassword","role":"admin","token":""}`,
		},
		{
			name: "admin get by Login",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?login=cmd@cmd.com",
			wantErr:    false,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() {
				mockUserService.On("GetUserByLogin", mock.AnythingOfType(ginContext), mock.AnythingOfType("string")).
					Return(domain.NewUser(domain.NewUserData{
						ID:        4,
						Login:     loginRegular,
						Password:  "hashpassword",
						Role:      "regular",
						Token:     "",
						CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
						UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
					}), nil).
					Once()
			},
			expectStatusCode:   http.StatusOK,
			expectResponseBody: `{"id":4,"login":"cmd@cmd.com","password":"hashpassword","role":"regular","token":""}`,
		},
		{
			name: "service error by id",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?id=4",
			wantErr:    true,
			mockGetUserByID: func() {
				mockUserService.On("GetUserByID", mock.AnythingOfType(ginContext), mock.AnythingOfType("int")).
					Return(domain.User{ // empty value
					},
						domain.ErrFailure).
					Once()
			},
			mockGetUserByLogin: func() { // empty value
			},
			expectStatusCode:   http.StatusInternalServerError,
			expectResponseBody: `{"error get users":"failure"}`,
		},
		{
			name: "service error by login",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?login=cmd@cmd.com",
			wantErr:    true,
			mockGetUserByID: func() { // empty value
			},
			mockGetUserByLogin: func() {
				mockUserService.On("GetUserByLogin", mock.AnythingOfType(ginContext), mock.AnythingOfType("string")).
					Return(domain.User{ // empty value
					},
						domain.ErrFailure).
					Once()
			},
			expectStatusCode:   http.StatusInternalServerError,
			expectResponseBody: `{"error get users":"failure"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user", tc.user)
			if tc.expectResponseBody == `{"error":"no user in context"}` {
				c.Set("user", nil)
			}

			url := fmt.Sprintf("/User%s", tc.queryInURL)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)
			c.Request.Header.Set("Content-Type", "application/json")
			tc.mockGetUserByID()
			tc.mockGetUserByLogin()
			h.GetUser(c)

			if tc.wantErr {
				require.Equal(t, tc.expectStatusCode, w.Code)
				require.Contains(t, w.Body.String(), tc.expectResponseBody)
			} else {
				assert.Equal(t, tc.expectStatusCode, w.Code)
				require.Equal(t, tc.expectResponseBody, w.Body.String())
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := mocks.NewIUserService(t)
	h := HTTPServer{
		userService: mockUserService,
	}
	testCases := []struct {
		name               string
		user               domain.User
		queryInURL         string
		wantErr            bool
		mockGetUsers       func()
		expectStatusCode   int
		expectResponseBody string
	}{
		{
			name: "OK",
			user: domain.NewUser(domain.NewUserData{
				ID:        1,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?limit=10&offset=0",
			wantErr:    false,
			mockGetUsers: func() {
				mockUserService.On("GetUsers", mock.AnythingOfType(ginContext), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return([]domain.User{
						domain.NewUser(domain.NewUserData{ID: 1, Login: "cmd@cmd.ru", Password: "hashpassword", Role: "admin", Token: "", CreatedAt: time.Time{ // empty value
						}, UpdatedAt:                                                                                                                           time.Time{ // empty value
						}}),
						domain.NewUser(domain.NewUserData{ID: 2, Login: "cmd@cmd.com", Password: "hashpassword", Role: "regular", Token: "", CreatedAt: time.Time{ // empty value
						}, UpdatedAt:                                                                                                                              time.Time{ // empty value
						}}),
						domain.NewUser(domain.NewUserData{ID: 3, Login: "cmd@cmd.org", Password: "hashpassword", Role: "regular", Token: "", CreatedAt: time.Time{ // empty value
						}, UpdatedAt:                                                                                                                              time.Time{ // empty value
						}}),
					}, nil).
					Once()
			},
			expectStatusCode:   http.StatusOK,
			expectResponseBody: `[{"id":1,"login":"cmd@cmd.ru","password":"hashpassword","role":"admin","token":""},{"id":2,"login":"cmd@cmd.com","password":"hashpassword","role":"regular","token":""},{"id":3,"login":"cmd@cmd.org","password":"hashpassword","role":"regular","token":""}]`,
		},
		{
			name: "not user in context",
			user: domain.User{ // empty value
			},
			queryInURL: "?limit=10&offset=0",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `{"error":"no user in context"}`,
		},
		{
			name: "not admin",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "regular",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?limit=10&offset=0",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusUnauthorized,
			expectResponseBody: `{"invalid user id or role":"access denied"}`,
		},
		{
			name:       "error get query limit",
			user:       domain.NewUser(domain.NewUserData{Role: "admin"}),
			queryInURL: "?limit=abc&offset=0",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `limit`,
		},
		{
			name:       "error get query offset",
			user:       domain.NewUser(domain.NewUserData{Role: "admin"}),
			queryInURL: "?limit=10&offset=abc",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `offset`,
		},
		{
			name:       "invalid limit",
			user:       domain.NewUser(domain.NewUserData{Role: "admin"}),
			queryInURL: "?limit=0&offset=0",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `limit-must-be-greater-then-zero`,
		},
		{
			name:       "invalid offset",
			user:       domain.NewUser(domain.NewUserData{Role: "admin"}),
			queryInURL: "?limit=10&offset=-1",
			wantErr:    true,
			mockGetUsers: func() { // empty value
			},
			expectStatusCode:   http.StatusBadRequest,
			expectResponseBody: `offset-must-be-greater-or-equal-then-zero`,
		},
		{
			name: "service error",
			user: domain.NewUser(domain.NewUserData{
				ID:        3,
				Login:     loginAdmin,
				Password:  "hashpassword",
				Role:      "admin",
				Token:     "",
				CreatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.June, 12, 15, 23, 13, 0, time.UTC),
			}),
			queryInURL: "?limit=10&offset=0",
			wantErr:    true,
			mockGetUsers: func() {
				mockUserService.On("GetUsers", mock.AnythingOfType(ginContext), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return([]domain.User{ // empty value
					},
						domain.ErrFailure).
					Once()
			},
			expectStatusCode:   http.StatusInternalServerError,
			expectResponseBody: `{"error get users":"failure"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user", tc.user)
			if tc.expectResponseBody == `{"error":"no user in context"}` {
				c.Set("user", nil)
			}

			url := fmt.Sprintf("/User%s", tc.queryInURL)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)
			c.Request.Header.Set("Content-Type", "application/json")
			tc.mockGetUsers()
			h.GetUsers(c)

			if tc.wantErr {
				require.Equal(t, tc.expectStatusCode, w.Code)
				require.Contains(t, w.Body.String(), tc.expectResponseBody)
			} else {
				require.Equal(t, tc.expectStatusCode, w.Code)
				require.Equal(t, tc.expectResponseBody, w.Body.String())
			}
		})
	}
}
