package config

import (
	"github.com/spf13/viper"
	"os"
)

type MongoSetting struct {
	URI        string `mapstructure:"uri" yaml:"uri"`
	DBNAME     string `mapstructure:"database" yaml:"database"`
	Timeout    int    `mapstructure:"timeout" yaml:"timeout"`
	Collection struct {
		Name string `mapstructure:"collection_name" yaml:"collection_name"`
	} `mapstructure:"collection" yaml:"collection"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret" yaml:"secret"`
	Expiry string `mapstructure:"expiry" yaml:"expiry"`
}

type ServerConfig struct {
	Port string `mapstructure:"port" yaml:"port"`
}

type Config struct {
	Mongo  MongoSetting `mapstructure:"mongo" yaml:"mongo"`
	JWT    JWTConfig    `mapstructure:"jwt" yaml:"jwt"`
	Server ServerConfig `mapstructure:"server" yaml:"server"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Mongo.URI = os.Getenv("MONGO_URI")
	config.JWT.Secret = os.Getenv("JWT_SECRET")
	config.JWT.Expiry = os.Getenv("JWT_EXPIRY")
	config.Server.Port = os.Getenv("SERVER_PORT")

	return &config, nil
}
