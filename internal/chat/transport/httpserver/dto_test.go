package httpserver

import (
	"net/http/httptest"
	"testing"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var userCtxKey = ctxKey("user")
	var productCtxKey = ctxKey("product")

	tests := []struct {
		name    string
		ctx     func() *gin.Context
		want    domain.User
		wantErr bool
	}{
		{
			name: "valid",
			ctx: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Set(userCtxKey.String(), domain.User{})
				return c
			},
			want:    domain.User{},
			wantErr: false,
		},
		{
			name: "invalid key",
			ctx: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Set(productCtxKey.String(), domain.User{})
				return c
			},
			want:    domain.User{},
			wantErr: true,
		},
		{
			name: "invalid struct",
			ctx: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Set(userCtxKey.String(), domain.NewUserData{})
				return c
			},
			want:    domain.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			c := tt.ctx()
			user, err := getUserFromContext(c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, user)
				assert.NoError(t, err)
			}
		})
	}
}
