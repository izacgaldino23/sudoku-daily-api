package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"sudoku-daily-api/src/domain"

	"golang.org/x/crypto/argon2"
)

type (
	passwordHasher struct {
		iterations  uint32
		memory      uint32
		parallelism uint8
		keyLen      uint32
		saltLen     uint32
	}
)

func NewPasswordHasher(
	iterations uint32,
	memory uint32,
	parallelism uint8,
	keyLen uint32,
	saltLen uint32,
) domain.PasswordHasher {
	return &passwordHasher{
		iterations:  iterations,
		memory:      memory,
		parallelism: parallelism,
		keyLen:      keyLen,
		saltLen:     saltLen,
	}
}

func (p *passwordHasher) Compare(password, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return err
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	hash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, comparisonHash) == 0 {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func (p *passwordHasher) Hash(password string) (string, error) {
	salt := make([]byte, p.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLen)

	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}
