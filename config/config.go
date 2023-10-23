package config

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KafkaBrokerHost string `envconfig:"KAFKA_BROKER_HOST" required:"true"`
	KafkaTopic      string `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaGroupId    string `envconfig:"KAFKA_GROUP_ID" required:"true"`
	MongodbDatabase string `envconfig:"MONGODB_DATABASE" required:"true"`
	MongodbHostName string `envconfig:"MONGODB_HOST_NAME" required:"true"`
	MongodbPort     int    `envconfig:"MONGODB_PORT" required:"true"`
}

var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

func Read(envFile string) (*Config, error) {
	if err := godotenvLoad(envFile); err != nil {
		return nil, errors.Unwrap(err)
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Unwrap(err)
	}
	return config, nil
}
