package config

import "github.com/spf13/viper"

type Config struct {
	ApiPort  string
	Database Database
}

type Database struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

func viperInit() error {
	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	return viper.ReadInConfig()
}

func Load() (*Config, error) {
	err := viperInit()
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}