package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/config"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account disabled")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtConfig config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtConfig: jwtConfig,
	}
}

type JWTClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	SchoolID *uuid.UUID      `json:"school_id,omitempty"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", ErrAccountDisabled
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	var birthDateStr *string
	if user.BirthDate != nil {
		str := user.BirthDate.Format("2006-01-02")
		birthDateStr = &str
	}

	return &domain.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwtConfig.AccessExpiry.Seconds()),
		User: domain.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			FullName:  user.FullName,
			PhotoURL:  user.PhotoURL,
			NUPTK:     user.NUPTK,
			NIP:       user.NIP,
			Gender:    user.Gender,
			BirthDate: birthDateStr,
			GTKType:   user.GTKType,
			Position:  user.Position,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, string, error) {
	tokenHash := hashToken(refreshToken)
	token, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, "", err
	}
	if token == nil {
		return nil, "", ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		s.tokenRepo.Delete(ctx, token.ID)
		return nil, "", ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, "", err
	}
	if user == nil || !user.IsActive {
		return nil, "", ErrAccountDisabled
	}

	// Delete old token
	s.tokenRepo.Delete(ctx, token.ID)

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", err
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return &domain.RefreshResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwtConfig.AccessExpiry.Seconds()),
	}, newRefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	token, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if token != nil {
		return s.tokenRepo.Delete(ctx, token.ID)
	}
	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.tokenRepo.DeleteByUserID(ctx, userID)
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) generateAccessToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		SchoolID: user.SchoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtConfig.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *AuthService) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	tokenID := uuid.New()
	rawToken := tokenID.String() + uuid.New().String()
	tokenHash := hashToken(rawToken)

	refreshToken := &domain.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.jwtConfig.RefreshExpiry),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	return rawToken, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
