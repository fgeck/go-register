package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fgeck/go-register/internal/repository"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrWeakPassword       = errors.New("password does not meet requirements")
)

type AuthService struct {
	repo          repository.Querier
	sessionExpiry time.Duration
	argon2Params  *argon2Params
	pepper        string // Application-wide secret for extra security
}

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewAuthService(repo repository.Querier) *AuthService {
	return &AuthService{
		repo:          repo,
		sessionExpiry: 24 * time.Hour * 7,      // 1 week sessions
		pepper:        "your-app-pepper-value", // In prod, load from secure config
		argon2Params: &argon2Params{
			memory:      64 * 1024, // 64MB
			iterations:  3,
			parallelism: 2,
			saltLength:  16, // 16 bytes
			keyLength:   32, // 32 bytes
		},
	}
}

// Register creates a new user with hashed password
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	// Validate password strength
	if len(password) < 8 {
		return ErrWeakPassword
	}

	// Check if user exists
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrUserExists
	}

	// Generate and store hashed password
	hash, err := s.generateHash(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.repo.CreateUser(ctx, repository.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Login verifies credentials and creates a session
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Simulate hashing to prevent timing attacks
		s.simulateHash()
		return "", ErrInvalidCredentials
	}

	// Verify password
	match, err := s.comparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !match {
		return "", ErrInvalidCredentials
	}

	// Create session
	session, err := s.repo.CreateSession(ctx, repository.CreateSessionParams{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.sessionExpiry),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return session.ID.String(), nil
}

// Logout invalidates a session
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// GetUserFromSession returns user if session is valid
func (s *AuthService) GetUserFromSession(ctx context.Context, sessionID string) (*repository.User, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.repo.DeleteSession(ctx, sessionID) // Cleanup expired session
		return nil, errors.New("session expired")
	}

	return s.repo.GetUser(ctx, session.UserID)
}

// --- Password Hashing Utilities ---

func (s *AuthService) generateHash(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, s.argon2Params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Combine password with pepper before hashing
	pepperedPassword := password + s.pepper

	// Generate hash
	hash := argon2.IDKey(
		[]byte(pepperedPassword),
		salt,
		s.argon2Params.iterations,
		s.argon2Params.memory,
		s.argon2Params.parallelism,
		s.argon2Params.keyLength,
	)

	// Encode to string format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		s.argon2Params.memory,
		s.argon2Params.iterations,
		s.argon2Params.parallelism,
		b64Salt,
		b64Hash,
	), nil
}

func (s *AuthService) comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Verify algorithm
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("unsupported algorithm")
	}

	// Parse parameters
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, fmt.Errorf("incompatible version")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	// Decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Hash the provided password with the same parameters
	pepperedPassword := password + s.pepper
	computedHash := argon2.IDKey(
		[]byte(pepperedPassword),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(storedHash)),
	)

	// Constant-time comparison
	return subtle.ConstantTimeCompare(computedHash, storedHash) == 1, nil
}

// simulateHash prevents timing attacks by running a dummy hash operation
func (s *AuthService) simulateHash() {
	argon2.IDKey(
		[]byte("dummy"),
		make([]byte, s.argon2Params.saltLength),
		s.argon2Params.iterations,
		s.argon2Params.memory,
		s.argon2Params.parallelism,
		s.argon2Params.keyLength,
	)
}
