package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	AWS AWSConfig
}

type AWSConfig struct {
	Region   string
	ID       string
	Secret   string
	Token    string
	Endpoint string
	Bucket   string
}

func NewConfig() Config {
	viper.AutomaticEnv()

	// NOTE: default values set for development with localstack
	viper.SetDefault("AWS_REGION", "eu-west-2")
	viper.SetDefault("AWS_ID", "test")
	viper.SetDefault("AWS_SECRET", "test")
	viper.SetDefault("AWS_TOKEN", "")
	viper.SetDefault("AWS_ENDPOINT", "http://localhost:4566")
	viper.SetDefault("AWS_S3_BUCKET", "fpc-test-bucket")

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Info().Err(err).Msg("config.Config failed to read file")
	}

	return Config{
		AWS: AWSConfig{
			Region:   viper.GetString("AWS_REGION"),
			ID:       viper.GetString("AWS_ID"),
			Secret:   viper.GetString("AWS_SECRET"),
			Token:    viper.GetString("AWS_TOKEN"),
			Endpoint: viper.GetString("AWS_ENDPOINT"),
			Bucket:   viper.GetString("AWS_S3_BUCKET"),
		},
	}
}
