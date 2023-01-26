package config

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

func NewConfig() (Config, error) {

	return Config{
		AWS: AWSConfig{
			Region:   "us-east-1",
			ID:       "test",
			Secret:   "test",
			Token:    "",
			Endpoint: "http://localhost:4566",
			Bucket:   "test-bucket",
		},
	}, nil
}
