package service

import (
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/domain/models"
	"compConfigurator/internal/repo"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     *repo.UsersRepo
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(users *repo.UsersRepo, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    jwtTTL,
	}
}

type TokenClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, email, password string, nickname *string) (models.User, string, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return models.User{}, "", domainErr.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", err
	}

	u, err := s.users.CreateUser(ctx, email, string(hash), nickname)
	if err != nil {
		return models.User{}, "", err
	}

	token, err := s.issueToken(u.UserID, u.RoleCode)
	if err != nil {
		return models.User{}, "", err
	}

	return u, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (models.User, string, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return models.User{}, "", domainErr.ErrInvalidInput
	}

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainErr.ErrNotFound) {
			return models.User{}, "", domainErr.ErrUnauthorized
		}
		return models.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return models.User{}, "", domainErr.ErrUnauthorized
	}

	token, err := s.issueToken(u.UserID, u.RoleCode)
	if err != nil {
		return models.User{}, "", err
	}

	return u, token, nil
}

func (s *AuthService) issueToken(userID int64, role string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, domainErr.ErrUnauthorized
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, domainErr.ErrUnauthorized
	}

	return claims, nil
}
