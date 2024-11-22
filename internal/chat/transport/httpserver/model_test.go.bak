package httpserver

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUserRequestValidate(t *testing.T) {
	cases := []struct {
		name          string
		in            UserRequest
		wantErr       bool
		expectedField []string
		expectedTag   []string
	}{
		{
			name: "valid user data",
			in: UserRequest{
				Login:    loginAdmin,
				Password: "12345678901234567890123456789012",
			},
			wantErr: false,
		},
		{
			name: "empty login",
			in: UserRequest{
				Login:    "",
				Password: "12345678901234567890123456789012",
			},
			wantErr:       true,
			expectedField: []string{"Login"},
			expectedTag:   []string{"required"},
		},
		{
			name: "wrong login format",
			in: UserRequest{
				Login:    "cmd@cmdru",
				Password: "12345678901234567890123456789012",
			},
			wantErr:       true,
			expectedField: []string{"Login"},
			expectedTag:   []string{"email"},
		},
		{
			name: "empty password",
			in: UserRequest{
				Login:    loginAdmin,
				Password: "",
			},
			wantErr:       true,
			expectedField: []string{"Password"},
			expectedTag:   []string{"required"},
		},
		{
			name: "short password error",
			in: UserRequest{
				Login:    loginAdmin,
				Password: "123",
			},
			wantErr:       true,
			expectedField: []string{"Password"},
			expectedTag:   []string{"min"},
		},
		{
			name: "long password error",
			in: UserRequest{
				Password: "123456789012345678901234567890123",
				Login:    loginAdmin,
			},
			wantErr:       true,
			expectedField: []string{"Password"},
			expectedTag:   []string{"max"},
		},
		{
			name: "long password error and wrong login",
			in: UserRequest{
				Login:    "cmd@cmdru",
				Password: "123456789012345678901234567890123",
			},
			wantErr:       true,
			expectedField: []string{"Login", "Password"},
			expectedTag:   []string{"email", "max"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.in.Validate()
			if c.wantErr {
				assert.Error(t, err)
				// Проверяем конкретные ошибки валидации
				validationErrors, ok := err.(validator.ValidationErrors)
				assert.True(t, ok)

				// Проверяем поле и тип ошибки
				if len(validationErrors) > 0 {
					for i := 0; i < len(validationErrors); i++ {
						ve := validationErrors[i]
						assert.Equal(t, c.expectedField[i], ve.Field())
						assert.Equal(t, c.expectedTag[i], ve.Tag())
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
