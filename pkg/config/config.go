package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	configEnv *Config
	loadMutex sync.Mutex
)

type Config struct {
	ApiPort  string   `mapstructure:"API_PORT"`
	Database Database `mapstructure:"DATABASE"`
	Auth     Auth     `mapstructure:"AUTH"`
	Limits   Limits   `mapstructure:"LIMITS"`

	Debug bool `mapstructure:"DEBUG"`
}

type Limits struct {
	MaxRequestCountGlobal int `mapstructure:"MAX_REQUEST_COUNT_GLOBAL"`
	MaxRequestCountUser   int `mapstructure:"MAX_REQUEST_COUNT_USER"`
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

	OidcEnabled  bool   `mapstructure:"OIDC_ENABLED"`
	OidcAudience string `mapstructure:"OIDC_AUDIENCE"`
}

type Database struct {
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`

	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
	SSLMode  string `mapstructure:"SSL_MODE"`

	MaxOpenConns int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int           `mapstructure:"MAX_IDLE_CONNS"`
	MaxLifetime  time.Duration `mapstructure:"MAX_LIFETIME"`
	Timeout      time.Duration `mapstructure:"TIMEOUT"`
}

func (d *Database) DSNPostgres() string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", d.Host, d.Username, d.Password, d.Name, d.Port, d.SSLMode)
}

func bindEnvs(v *viper.Viper, c reflect.Type, prefix []string) error {
	tag := "mapstructure"
	for i := range c.NumField() {
		field := c.Field(i)

		if tagValue, ok := field.Tag.Lookup(tag); ok {
			if field.Type.Kind() == reflect.Struct {
				err := bindEnvs(v, field.Type, append(prefix, tagValue))
				if err != nil {
					return err
				}
			} else {
				current := append(prefix, tagValue)

				err := v.BindEnv(strings.Join(current, "."), strings.Join(current, "_"))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func viperInit() error {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("env")

	if envPath := os.Getenv("ENV_FILE"); envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	v.AutomaticEnv()

	configEnv = &Config{}
	err := bindEnvs(v, reflect.TypeOf(*configEnv), []string{})
	if err != nil {
		return err
	}

	err = v.Unmarshal(configEnv)
	if err != nil {
		return err
	}

	if os.Getenv("ENV") == "local" {
		configEnv.ApiPort = "127.0.0.1:" + configEnv.ApiPort
	} else {
		configEnv.ApiPort = "0.0.0.0:" + configEnv.ApiPort
	}

	return nil
}

func Load() error {
	loadMutex.Lock()
	defer loadMutex.Unlock()

	if configEnv != nil {
		return nil
	}

	return viperInit()
}

func GetConfig() *Config {
	if configEnv == nil {
		_ = Load()
	}

	return configEnv
}
