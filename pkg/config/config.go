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
}

type Database struct {
	Host         string `mapstructure:"HOST"`
	Port         string `mapstructure:"PORT"`
	Username     string `mapstructure:"USERNAME"`
	Password     string `mapstructure:"PASSWORD"`
	Name         string `mapstructure:"NAME"`
	SSLMode      string `mapstructure:"SSL_MODE"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxLifetime  int    `mapstructure:"MAX_LIFETIME"`
}

func (d *Database) DSNPostgres() string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", d.Host, d.Username, d.Password, d.Name, d.Port, d.SSLMode)
}

func viperInit() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	v.SetConfigType("env")

	name := os.Getenv("ENV")
	if name == "" {
		name = "local"
	}
	// evita que alguém passe "local.env" na ENV
	name = strings.TrimSuffix(name, ".env")

	v.SetConfigName(name) // sem a extensão
	v.AddConfigPath(".")  // procura no cwd

	if _, err := os.Stat(name + ".env"); err != nil {
		return nil, err
	}

	v.SetDefault("DATABASE.SSL_MODE", "disable")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, err
		}
	}

	for _, key := range v.AllKeys() {
		if strings.Contains(key, "_") {
			newKey := strings.Replace(key, "_", ".", -1)
			v.Set(newKey, v.Get(key))
		}
	}

	c := &Config{}

	err := v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Load() (err error) {
	if configEnv == nil {
		once := sync.Once{}
		once.Do(func() {
			configEnv, err = viperInit()
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
