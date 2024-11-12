package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
)

// TODO: move to secrets
var jwtSecretKey = []byte(config.JwtKey)

// TokenService is a token service
type TokenService struct {
	repo IUserRepository
	ttl  time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(
	repo IUserRepository,
	ttl time.Duration,
) TokenService {
	return TokenService{
		repo: repo,
		ttl:  ttl,
	}
}

// UserClaims ...
type UserClaims struct {
	AuthID    int    `json:"auth_id"`
	AuthLogin string `json:"auth_login"`
	AuthRole  string `json:"auth_role"`
	jwt.StandardClaims
}

// GenerateToken generates a token
// func (s TokenService) GenerateToken(user domain.User) (string, error) {
func (s TokenService) GenerateToken(ctx context.Context, login, password string) (string, error) {
	domainUser, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("invaldRequest: %v", err.Error())
	}

	if !utils.CheckPasswordHash(password, domainUser.Password()) {
		return "", fmt.Errorf("error: invalid-password")
	}
	fmt.Printf("func GenerateToken: domainUser: %+v\n", domainUser)
	payload := UserClaims{
		AuthID:    domainUser.ID(),
		AuthLogin: domainUser.Login(),
		AuthRole:  domainUser.Role(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			// ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			ExpiresAt: time.Now().Add(config.Cfg.TokenTimeDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return t, nil
}

// GetUser ...
func (s TokenService) GetUser(token string) (domain.User, error) {
	var userClaims UserClaims
	t, err := jwt.ParseWithClaims(token, &userClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to parse a token: %w", err)
	}
	if !t.Valid {
		return domain.User{}, errors.New("invalid token")
	}
	user := userClaimsToDomainUser(userClaims)
	return user, nil
}

func userClaimsToDomainUser(claims UserClaims) domain.User {
	return domain.NewUser(domain.NewUserData{
		ID:    claims.AuthID,
		Login: claims.AuthLogin,
		Role:  claims.AuthRole,
	})
}
