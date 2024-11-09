package config

import (
	"OzonHW1/server/internal/infra/kafka"
)

type ProducerConfig struct {
	Topic string
}

type Config struct {
	Kafka    kafka.Config
	Producer ProducerConfig
}

func NewConfig(topic, kafkaBroker string) Config {
	return Config{

		Kafka: kafka.Config{
			Brokers: []string{
				kafkaBroker,
			},
		},
		Producer: ProducerConfig{
			Topic: topic,
		},
	}
}
