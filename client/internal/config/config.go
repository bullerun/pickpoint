package config

import "OzonHW1/client/internal/infra/kafka"

type Config struct {
	KafkaConfig kafka.Config
}

func NewConfig(bootstrap string) Config {
	return Config{
		KafkaConfig: kafka.Config{
			Brokers: []string{
				bootstrap,
			},
		},
	}
}
