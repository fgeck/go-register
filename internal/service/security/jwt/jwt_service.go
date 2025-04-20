package jwt

import "github.com/google/uuid"

type JwtServiceInterface interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(token string) (string, error)
	ParseToken(token string) (string, error)
	ExtractClaims(token string) (map[string]interface{}, error)
	ExtractUserId(token string) (string, error)
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
	// Implementation for generating JWT token
	return "", nil
}
func (s *JwtService) ValidateToken(token string) (string, error) {
	// Implementation for validating JWT token
	return "", nil
}
func (s *JwtService) ParseToken(token string) (string, error) {
	// Implementation for parsing JWT token
	return "", nil
}
func (s *JwtService) ExtractClaims(token string) (map[string]interface{}, error) {
	// Implementation for extracting claims from JWT token
	return nil, nil
}
func (s *JwtService) ExtractUserId(token string) (string, error) {
	// Implementation for extracting user ID from JWT token
	return "", nil
}
