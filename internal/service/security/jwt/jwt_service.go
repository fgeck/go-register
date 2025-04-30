package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/google/uuid"
)

type JwtServiceInterface interface {
	GenerateToken(user *user.UserDto) (string, error)
	ValidateAndExtractClaims(givenToken string) (*Claims, error)
}

type JwtService struct {
	secretKey  string
	issuer     string
	expiration int64
}

func NewJwtService(secretKey, issuer string, expiration int64) *JwtService {
	return &JwtService{
		secretKey:  secretKey,
		issuer:     issuer,
		expiration: expiration,
	}
}

func (s *JwtService) GenerateToken(user *user.UserDto) (string, error) {
	if user.ID == uuid.Nil || user.ID.String() == "" {
		return "", fmt.Errorf("userID is empty")
	}
	if user.Role.Name == "" {
		return "", fmt.Errorf("userRole is empty")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		USER_ID:   user.ID.String(),
		USER_ROLE: user.Role.Name,
		"iss":     s.issuer,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Duration(s.expiration) * time.Second).Unix(),
		"nbf":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) ValidateAndExtractClaims(givenToken string) (*Claims, error) {
	token, err := jwt.Parse(givenToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return s.buildClaims(claims)
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (s *JwtService) buildClaims(jwtClaims jwt.MapClaims) (*Claims, error) {
	userId, ok := jwtClaims[USER_ID]
	if !ok {
		return nil, fmt.Errorf("missing userId claim")
	}
	userRole, ok := jwtClaims[USER_ROLE]
	if !ok {
		return nil, fmt.Errorf("missing userRole claim")
	}
	if userId == nil || userRole == nil {
		return nil, fmt.Errorf("userId or userRole claim is nil. claims: %v", jwtClaims)

	}

	return NewClaims(userId.(string), userRole.(string)), nil
}
