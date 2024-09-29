package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "valid test",
			password:  "123456",
			wantError: false,
		},
		{
			name:      "password is too long",
			password:  "1234567890123456789012345678901234567890123456789012345678901234567890123",
			wantError: true,
		},
	}
	for _, tc := range testCases {
		hash, err := hashPassword(tc.password)
		if tc.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tc.password, hash)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "123456"

	hash, err := hashPassword(password)
	assert.NoError(t, err)

	testCases := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash,
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}
