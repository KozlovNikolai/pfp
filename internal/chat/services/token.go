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
	AuthID       int    `json:"auth_id"`
	AuthAccount  string `json:"account"`
	AuthLogin    string `json:"auth_login"`
	AuthUserType string `json:"auth_user_type"`
	jwt.StandardClaims
}

// GenerateToken generates a token
func (s TokenService) GenerateToken(ctx context.Context, account, login, password string) (domain.User, string, error) {
	domainUserChat, err := s.repo.GetUserByLogin(ctx, account, login)
	if err != nil {
		return domain.User{}, "", fmt.Errorf("invaldRequest: %v", err.Error())
	}

	if !utils.CheckPasswordHash(password, domainUserChat.Password()) {
		return domain.User{}, "", fmt.Errorf("error: invalid-password")
	}
	fmt.Printf("func GenerateToken: domainUser: %+v\n", domainUserChat)
	payload := UserClaims{
		AuthID:       domainUserChat.ID(),
		AuthLogin:    domainUserChat.Login(),
		AuthUserType: domainUserChat.UserType(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			// ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			ExpiresAt: time.Now().Add(config.Cfg.TokenTimeDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return domain.User{}, "", fmt.Errorf("failed to sign token: %w", err)
	}

	return domainUserChat, t, nil
}

// GetUser ...
func (s TokenService) GetUser(ctx context.Context, token string) (domain.User, error) {
	_ = ctx
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
		ID:       claims.AuthID,
		Account:  claims.AuthAccount,
		Login:    claims.AuthLogin,
		UserType: claims.AuthUserType,
	})
}

func (s TokenService) GenerateTokenForRegisteredUsers(ctx context.Context, user domain.User) (string, error) {

	payload := UserClaims{
		AuthID:       user.ID(),
		AuthAccount:  user.Account(),
		AuthLogin:    user.Login(),
		AuthUserType: user.UserType(),
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
