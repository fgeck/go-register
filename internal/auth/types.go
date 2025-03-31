package auth

import "time"

type AuthConfig struct {
	SessionExpiry time.Duration
	Pepper        string
	Argon2Params  *Argon2Params
}

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}
