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

func NewConfig() *Config {
	viper.AutomaticEnv()

	// set defaults to empty strings
	viper.SetDefault("AWS_REGION", "")
	viper.SetDefault("AWS_ID", "")
	viper.SetDefault("AWS_SECRET", "")
	viper.SetDefault("AWS_TOKEN", "")
	viper.SetDefault("AWS_ENDPOINT", "")
	viper.SetDefault("AWS_S3_BUCKET", "")

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Info().Err(err).Msg("config.Config failed to read file")
	}

	return &Config{
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
