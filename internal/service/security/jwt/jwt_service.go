package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JwtServiceInterface interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(token string) (string, error)
	ParseToken(token string) (string, error)
	ExtractClaims(token string) (map[string]interface{}, error)
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

func (s *JwtService) GenerateToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId.String(),
		"issuer": s.issuer,
		"exp":    time.Now().Add(time.Duration(s.expiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) ValidateToken(token string) (string, error) {
	return "", nil
}

func (s *JwtService) ParseToken(token string) (string, error) {
	return "", nil
}

func (s *JwtService) ExtractClaims(token string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

