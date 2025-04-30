package jwt_test

import (
	"testing"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/fgeck/go-register/internal/service/security/jwt"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_PRIVEATE_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC0C34XotpaxNWx
YG6s64stT40HjCt4zo4naCc9UPmL6+FZ9ZTcVCJOdzI8/WoyFjEni0CPBlvNUoAi
3W+sM9g+6m79VP6RrZ18+75WsDfz2pr0VLAVOQNog3Q+WWWD4vV+J0C4PnL8xvKA
bMJx3EdCscDApPFq0QCA9+RVCr9mOlVe5nSoVBnKY9SfirSbptgJqCYO8pE8rtRQ
9Xe9q0LxQtfTKxoZD4M8kqcOl7d1lPhgQbOTdarIge+4uacg4yhsSD6UVnOKS8zW
e7hIORoby4h92Lo3W2zTh0KMWmzrbXASMIP1WqRF0kVpW/E917sf5eDi2023OrCe
qgMbwSBTAgMBAAECggEAJDgrpZWdV6VBV+2OVjsMRJE8Tchk9miXMFZDjpI7oWpS
a0Z8K9bBEAfqk1pngqv0N4BL/HnK/gMaw+jIDlxfpEiFC3GNxMCobfw2zjmlB+ly
QrTGt35AsUXAnMpfIakGudorquTlPPTI1A0NENq/eytHG3oTFun0r/0uce03k9jR
nxgMu8qPBFBdeTh2QpiCFrZqsanPKJRJiYopRLshTLaRdH0uwutpYXJg7vcgJd7+
8TSx5Z7xnAep5cqbeUxvoSsiucyHN8FDZpo4zlQ4UzSVfPC+lnywRIFupWfjW1PB
96/1tO2Ea3aRVgSmzEd0DYWHkUgyIJnqGQ0qpXpIAQKBgQDqMl9JTlfa7+OtAQ5Q
EfXElWFHvwRlv2+yHm71KshvNwIdg5POM8MCGQdNgRYSdcar7KsPyDy+rdKV18UC
NwCwWsv88z3LwaRjGc6Oj6658CZuHuQO4Sm/lifaeLY+nVhw0jhmg2DfDiJsgdsF
mGLOO8+t/wCg1XIVrX8GJqhuGQKBgQDEzoPVN/rK/e4Zw1duk1jAhJPjv2B76KzD
K+3pbiPz4syRtTWuEYdLc8gyYYtK9eF3oc8fjgCMgzQ09Q/1nU3KceXP0yFVrmRm
ejy1XFORHRpfAGNbqYxDleIkNz4bmg91Ff1yef2gsZIVjZnWMPxgDna4qSPt9DDJ
LPzvi/23SwKBgQDF9iBPUba3rRERuxPDIPtS2UYqpE9uRjx/HnSCLlDQmXnjQsZc
hapwCoH+xH/IyN9PkjUimQqnzzxzRrkT3zRo3ccSIPX6VsvCrRzJqrByIYoKiXgT
D8b/WEiFxoWeNdh9PWVJWgI3abY1bCqb9yyF0U8Cb8uzJ9lQc6Asrd6veQKBgCHu
RPZmz1teCkXw0ssipkOS1/iFDzpttBN2KG99aL9sk75vUpDvPrc4gASHor9KwxOg
Fximn9uZ509WDOlYtIe5uVhqWy3tgivU2VCfWV0Een50j6zG/4LLfZCm4ZNarV2P
bAHnnF2vH7ONlT9DdM+OztMpfiNRXXPhyL34EccfAoGBAJXRCPoBpUWmTkF/5HhN
3Y9sGakouw0HP5eNq1btZ51yytzSWu3/rL1JifmZRnCkOzLu+LYsj+wi7tmJVSpy
U5IfQg0LmP5ruurcoJ6V8MSifkrPan9m3Uz+S1ezJNr/XV4T5QTD04zCnE6yoJMr
+CarzOg1+2CoUNsIlJJ6SQmJ
-----END PRIVATE KEY-----
`
	TEST_SECRET = "SuperVal!d1@asdawe36"
)

func TestGenerateToken(t *testing.T) {
	jwtService := jwt.NewJwtService(TEST_SECRET, "test-issuer", 3600)

	t.Run("Valid Input", func(t *testing.T) {
		userDto := &user.UserDto{
			ID:   uuid.New(),
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Empty User ID", func(t *testing.T) {
		userDto := &user.UserDto{
			ID:   uuid.Nil,
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Empty User Role", func(t *testing.T) {
		userDto := &user.UserDto{
			ID:   uuid.New(),
			Role: user.UserRole{Name: ""},
		}

		token, err := jwtService.GenerateToken(userDto)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestValidateAndExtractClaims(t *testing.T) {
	jwtService := jwt.NewJwtService(TEST_SECRET, "test-issuer", 3600)

	t.Run("Valid Token", func(t *testing.T) {
		userID := uuid.New()
		userDto := &user.UserDto{
			ID:   userID,
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		assert.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(token)
		assert.NoError(t, err)
		assert.Equal(t, userID.String(), extractedClaims.UserId)
		assert.Equal(t, user.UserRoleAdmin.Name, extractedClaims.UserRole)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		token := "invalid-token"

		extractedUser, err := jwtService.ValidateAndExtractClaims(token)
		assert.Error(t, err)
		assert.Nil(t, extractedUser)
	})

	t.Run("No HMAC Signing Method Used", func(t *testing.T) {
		privateKey, err := jwtGo.ParseRSAPrivateKeyFromPEM([]byte(TEST_PRIVEATE_KEY))
		assert.NoError(t, err)

		// Create a token with RS256 signing method
		claims := jwtGo.MapClaims{
			"userId":   uuid.New().String(),
			"userRole": "admin",
			"iss":      "test-issuer",
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour).Unix(),
			"nbf":      time.Now().Unix(),
		}

		token := jwtGo.NewWithClaims(jwtGo.SigningMethodRS256, claims)
		signedToken, err := token.SignedString(privateKey)
		assert.NoError(t, err)

		// Validate the token using the JwtService
		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		assert.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("No User ID in Parsed Token", func(t *testing.T) {
		// Create a token without a USER_ID claim
		claims := jwtGo.MapClaims{
			"userRole": "admin",
			"iss":      "test-issuer",
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour).Unix(),
		}

		token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TEST_SECRET))
		assert.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		assert.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "missing userId claim")
	})

	t.Run("No User Role in Parsed Token", func(t *testing.T) {
		// Create a token without a USER_ROLE claim
		claims := jwtGo.MapClaims{
			"userId": uuid.New().String(),
			"iss":    "test-issuer",
			"iat":    time.Now().Unix(),
			"exp":    time.Now().Add(time.Hour).Unix(),
		}

		token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TEST_SECRET))
		assert.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		assert.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "missing userRole claim")
	})
}
