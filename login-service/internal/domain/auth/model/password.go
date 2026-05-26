package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHash is a value object representing a bcrypt-hashed password.
type PasswordHash string

// NewPasswordHash creates a PasswordHash from a raw password by hashing it.
func NewPasswordHash(rawPassword string) (PasswordHash, error) {
	if rawPassword == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	if len(rawPassword) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return PasswordHash(string(hash)), nil
}

// PasswordHashFromHash creates a PasswordHash from an existing hash string.
func PasswordHashFromHash(hash string) PasswordHash {
	return PasswordHash(hash)
}

// Verify checks a raw password against the stored hash.
func (h PasswordHash) Verify(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(string(h)), []byte(rawPassword))
	return err == nil
}

// String returns the hash string.
func (h PasswordHash) String() string {
	return string(h)
}
