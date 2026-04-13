package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetConfig() {
	configEnv = nil
}

type envMap map[string]string

func (e envMap) set() {
	for k, v := range e {
		os.Setenv(k, v)
	}
}

func (e envMap) unset() {
	for k := range e {
		os.Unsetenv(k)
	}
}

func setupEnv(t *testing.T, env envMap) func() {
	env.set()
	return func() { env.unset() }
}

func createEnvFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envFilePath, []byte(content), 0644)
	require.NoError(t, err)
	return envFilePath
}

var defaultEnvVars = envMap{
	"API_PORT":                    "8080",
	"DATABASE_HOST":               "localhost",
	"DATABASE_PORT":               "5432",
	"DATABASE_USERNAME":           "user",
	"DATABASE_PASSWORD":           "pass",
	"DATABASE_NAME":               "sudoku",
	"DATABASE_SSL_MODE":           "require",
	"DATABASE_MIGRATIONS_PATH":    "./migrations",
	"DATABASE_MAX_OPEN_CONNS":     "25",
	"DATABASE_MAX_IDLE_CONNS":     "5",
	"DATABASE_MAX_LIFETIME":       "300s",
	"LIMITS_TIMEOUT":              "3s",
	"DEBUG":                       "true",
	"AUTH_ITERATIONS":             "65536",
	"AUTH_MEMORY":                 "65536",
	"AUTH_PARALLELISM":            "4",
	"AUTH_KEY_LEN":                "32",
	"AUTH_SALT_LEN":               "16",
	"AUTH_SECRET_KEY":             "test-secret-key",
	"AUTH_ACCESS_TOKEN_DURATION":  "3600",
	"AUTH_REFRESH_TOKEN_DURATION": "86400",
	"AUTH_OIDC_ENABLED":           "true",
	"AUTH_OIDC_AUDIENCE":          "test-audience",
	"ENV":                         "production",
}

func TestLoadConfig(t *testing.T) {
	resetConfig()
	defer setupEnv(t, defaultEnvVars)()

	err := Load()
	require.NoError(t, err)

	cfg := GetConfig()
	assert.NotNil(t, cfg)

	assert.Equal(t, "0.0.0.0:8080", cfg.ApiPort)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "user", cfg.Database.Username)
	assert.Equal(t, "pass", cfg.Database.Password)
	assert.Equal(t, "sudoku", cfg.Database.Name)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "./migrations", cfg.Database.MigrationsPath)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, time.Second*300, cfg.Database.MaxLifetime)
	assert.True(t, cfg.Debug)
	assert.Equal(t, uint32(65536), cfg.Auth.Iterations)
	assert.Equal(t, uint32(65536), cfg.Auth.Memory)
	assert.Equal(t, uint8(4), cfg.Auth.Parallelism)
	assert.Equal(t, uint32(32), cfg.Auth.KeyLen)
	assert.Equal(t, uint32(16), cfg.Auth.SaltLen)
	assert.Equal(t, "test-secret-key", cfg.Auth.SecretKey)
	assert.Equal(t, 3600, cfg.Auth.AccessTokenDuration)
	assert.Equal(t, 86400, cfg.Auth.RefreshTokenDuration)
	assert.True(t, cfg.Auth.OidcEnabled)
	assert.Equal(t, "test-audience", cfg.Auth.OidcAudience)
}

func TestLoadConfigLocalEnv(t *testing.T) {
	resetConfig()
	env := envMap{
		"API_PORT":      "3000",
		"DATABASE_HOST": "localhost",
		"ENV":           "local",
	}
	defer setupEnv(t, env)()

	err := Load()
	require.NoError(t, err)

	cfg := GetConfig()
	assert.Equal(t, "127.0.0.1:3000", cfg.ApiPort)
}

func TestGetConfigWithoutLoad(t *testing.T) {
	resetConfig()
	env := envMap{
		"API_PORT":      "9000",
		"DATABASE_HOST": "testhost",
	}
	defer setupEnv(t, env)()

	cfg := GetConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, "testhost", cfg.Database.Host)
}

func TestDSNPostgres(t *testing.T) {
	resetConfig()
	env := envMap{
		"DATABASE_HOST":     "dbhost",
		"DATABASE_PORT":     "5432",
		"DATABASE_USERNAME": "user",
		"DATABASE_PASSWORD": "pass",
		"DATABASE_NAME":     "testdb",
		"DATABASE_SSL_MODE": "disable",
	}
	defer setupEnv(t, env)()

	err := Load()
	require.NoError(t, err)

	cfg := GetConfig()
	dsn := cfg.Database.DSNPostgres()
	assert.Equal(t, "host=dbhost user=user password=pass dbname=testdb port=5432 sslmode=disable", dsn)
}

func TestLoadConfigFromEnvFile(t *testing.T) {
	resetConfig()

	envFileContent := `API_PORT=9090
DATABASE_HOST=dbfromfile
DATABASE_PORT=5433
DATABASE_USERNAME=fileuser
DATABASE_PASSWORD=filepass
DATABASE_NAME=filedb
DATABASE_SSL_MODE=require
DATABASE_MIGRATIONS_PATH=/migrations
DATABASE_MAX_OPEN_CONNS=50
DATABASE_MAX_IDLE_CONNS=10
DATABASE_MAX_LIFETIME=600s
LIMITS_TIMEOUT=3s
DEBUG=false
AUTH_ITERATIONS=8192
AUTH_MEMORY=8192
AUTH_PARALLELISM=2
AUTH_KEY_LEN=64
AUTH_SALT_LEN=32
AUTH_SECRET_KEY=file-secret-key
AUTH_ACCESS_TOKEN_DURATION=1800
AUTH_REFRESH_TOKEN_DURATION=43200
AUTH_OIDC_ENABLED=false
AUTH_OIDC_AUDIENCE=file-audience
`

	envFilePath := createEnvFile(t, envFileContent)
	defer setupEnv(t, envMap{"ENV_FILE": envFilePath})()

	err := Load()
	require.NoError(t, err)

	cfg := GetConfig()
	assert.NotNil(t, cfg)

	assert.Equal(t, "0.0.0.0:9090", cfg.ApiPort)
	assert.Equal(t, "dbfromfile", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "fileuser", cfg.Database.Username)
	assert.Equal(t, "filepass", cfg.Database.Password)
	assert.Equal(t, "filedb", cfg.Database.Name)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "/migrations", cfg.Database.MigrationsPath)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, time.Second*600, cfg.Database.MaxLifetime)
	assert.False(t, cfg.Debug)
	assert.Equal(t, uint32(8192), cfg.Auth.Iterations)
	assert.Equal(t, uint32(8192), cfg.Auth.Memory)
	assert.Equal(t, uint8(2), cfg.Auth.Parallelism)
	assert.Equal(t, uint32(64), cfg.Auth.KeyLen)
	assert.Equal(t, uint32(32), cfg.Auth.SaltLen)
	assert.Equal(t, "file-secret-key", cfg.Auth.SecretKey)
	assert.Equal(t, 1800, cfg.Auth.AccessTokenDuration)
	assert.Equal(t, 43200, cfg.Auth.RefreshTokenDuration)
	assert.False(t, cfg.Auth.OidcEnabled)
	assert.Equal(t, "file-audience", cfg.Auth.OidcAudience)
}

func TestLoadConfigEnvVarsOverrideEnvFile(t *testing.T) {
	resetConfig()

	envFilePath := createEnvFile(t, `DATABASE_HOST=envfile-host
DATABASE_PORT=5434`)

	env := envMap{
		"ENV_FILE":      envFilePath,
		"DATABASE_HOST": "env-var-host",
		"DATABASE_PORT": "5432",
	}
	defer setupEnv(t, env)()

	err := Load()
	require.NoError(t, err)

	cfg := GetConfig()
	assert.Equal(t, "env-var-host", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
}

func TestLoadConfigInvalidEnvFile(t *testing.T) {
	resetConfig()
	defer setupEnv(t, envMap{"ENV_FILE": "/nonexistent/path/.env"})()

	err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load env file")
}
