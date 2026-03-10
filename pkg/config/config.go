package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	configEnv *Config
)

type Config struct {
	ApiPort  string   `mapstructure:"API_PORT"`
	Database Database `mapstructure:"DATABASE"`
	Auth     Auth     `mapstructure:"AUTH"`
}

type Auth struct {
	Iterations  uint32 `mapstructure:"ITERATIONS"`
	Memory      uint32 `mapstructure:"MEMORY"`
	Parallelism uint8  `mapstructure:"PARALLELISM"`
	KeyLen      uint32 `mapstructure:"KEY_LEN"`
	SaltLen     uint32 `mapstructure:"SALT_LEN"`

	SecretKey            string `mapstructure:"SECRET_KEY"`
	AccessTokenDuration  int    `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration int    `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type Database struct {
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`

	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
	SSLMode  string `mapstructure:"SSL_MODE"`

	MaxOpenConns int `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int `mapstructure:"MAX_IDLE_CONNS"`
	MaxLifetime  int `mapstructure:"MAX_LIFETIME"`
}

func (d *Database) DSNPostgres() string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", d.Host, d.Username, d.Password, d.Name, d.Port, d.SSLMode)
}

func viperInit() error {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("env")

	name := os.Getenv("ENV")
	if name == "" {
		name = "local"
	}
	name = strings.TrimSuffix(name, ".env")

	v.SetConfigName(name)
	v.AddConfigPath(".")

	v.SetDefault("DATABASE.SSL_MODE", "disable")

	if _, err := os.Stat(name + ".env"); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	v.AutomaticEnv()

	_ = v.BindEnv("DATABASE.HOST", "DATABASE_HOST")
	_ = v.BindEnv("DATABASE.PORT", "DATABASE_PORT")
	_ = v.BindEnv("DATABASE.USERNAME", "DATABASE_USERNAME")
	_ = v.BindEnv("DATABASE.PASSWORD", "DATABASE_PASSWORD")
	_ = v.BindEnv("DATABASE.NAME", "DATABASE_NAME")
	_ = v.BindEnv("DATABASE.SSL_MODE", "DATABASE_SSL_MODE")
	_ = v.BindEnv("DATABASE.MIGRATIONS_PATH", "MIGRATIONS_PATH")

	_ = v.BindEnv("API_PORT")
	_ = v.BindEnv("AUTH.ITERATIONS", "AUTH_ITERATIONS")
	_ = v.BindEnv("AUTH.MEMORY", "AUTH_MEMORY")
	_ = v.BindEnv("AUTH.PARALLELISM", "AUTH_PARALLELISM")
	_ = v.BindEnv("AUTH.KEY_LEN", "AUTH_KEY_LEN")
	_ = v.BindEnv("AUTH.SALT_LEN", "AUTH_SALT_LEN")
	_ = v.BindEnv("AUTH.SECRET_KEY", "AUTH_SECRET_KEY")
	_ = v.BindEnv("AUTH.ACCESS_TOKEN_DURATION", "AUTH_ACCESS_TOKEN_DURATION")
	_ = v.BindEnv("AUTH.REFRESH_TOKEN_DURATION", "AUTH_REFRESH_TOKEN_DURATION")

	configEnv = &Config{}

	err := v.Unmarshal(configEnv)
	if err != nil {
		return err
	}

	if name == "local" {
		configEnv.ApiPort = "127.0.0.1:" + configEnv.ApiPort
	} else {
		configEnv.ApiPort = "0.0.0.0:" + configEnv.ApiPort
	}

	return nil
}

func Load() (err error) {
	if configEnv == nil {
		once := sync.Once{}
		once.Do(func() {
			err = viperInit()
		})
	}

	return
}

func GetConfig() *Config {
	if configEnv == nil {
		_ = Load()
	}

	return configEnv
}
