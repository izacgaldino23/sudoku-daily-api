package services_test

import (
	"strings"
	"testing"

	"sudoku-daily-api/src/services"
)

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := services.NewPasswordHasher(1, 64*1024, 1, 32, 16)

	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "password123"},
		{"empty password", ""},
		{"special characters", "p@ssw0rd!#$%^&*()"},
		{"unicode password", "密码密码"},
		{"long password", strings.Repeat("a", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := hasher.Hash(tt.password)
			if err != nil {
				t.Fatalf("Hash() error = %v", err)
			}

			if hash == "" {
				t.Error("Hash() returned empty string")
			}

			if !strings.HasPrefix(hash, "$argon2id$") {
				t.Errorf("Hash() = %q; want to start with $argon2id$", hash)
			}

			parts := strings.Split(hash, "$")
			if len(parts) != 6 {
				t.Errorf("Hash() format = %d parts; want 6", len(parts))
			}
		})
	}
}

func TestPasswordHasher_Compare(t *testing.T) {
	hasher := services.NewPasswordHasher(1, 64*1024, 1, 32, 16)

	hash, err := hasher.Hash("correctpassword")
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name        string
		password    string
		encodedHash string
		wantErr     bool
	}{
		{
			name:        "correct password",
			password:    "correctpassword",
			encodedHash: hash,
			wantErr:     false,
		},
		{
			name:        "incorrect password",
			password:    "wrongpassword",
			encodedHash: hash,
			wantErr:     true,
		},
		{
			name:        "empty password against non-empty hash",
			password:    "",
			encodedHash: hash,
			wantErr:     true,
		},
		{
			name:        "invalid hash format - missing parts",
			password:    "password",
			encodedHash: "$argon2id$v=19$m=65536,t=1,p=1",
			wantErr:     true,
		},
		{
			name:        "invalid hash format - wrong prefix",
			password:    "password",
			encodedHash: "$bcrypt$v=19$m=65536,t=1,p=1$AAAA$BBBB",
			wantErr:     true,
		},
		{
			name:        "invalid base64 salt",
			password:    "password",
			encodedHash: "$argon2id$v=19$m=65536,t=1,p=1$!!!invalid$AAAA",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Compare(tt.password, tt.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v; wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordHasher_Compare_RejectsDifferentPasswords(t *testing.T) {
	hasher := services.NewPasswordHasher(1, 64*1024, 1, 32, 16)

	hash1, _ := hasher.Hash("password1")
	hash2, _ := hasher.Hash("password2")

	if err := hasher.Compare("password1", hash2); err == nil {
		t.Error("Compare() should reject different password hashes")
	}

	if err := hasher.Compare("password2", hash1); err == nil {
		t.Error("Compare() should reject different password hashes")
	}
}

func TestPasswordHasher_GenerateUniqueHashes(t *testing.T) {
	hasher := services.NewPasswordHasher(1, 64*1024, 1, 32, 16)

	password := "samepassword"

	hash1, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	hash2, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash1 == hash2 {
		t.Error("Hash() should generate unique hashes for same password (due to random salt)")
	}

	if err := hasher.Compare(password, hash1); err != nil {
		t.Errorf("Compare() should validate first hash: %v", err)
	}

	if err := hasher.Compare(password, hash2); err != nil {
		t.Errorf("Compare() should validate second hash: %v", err)
	}
}
