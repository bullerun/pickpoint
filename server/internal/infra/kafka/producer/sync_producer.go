package producer

import (
	"OzonHW1/server/internal/infra/kafka"
	"fmt"

	"github.com/IBM/sarama"
)

func NewSyncProducer(conf kafka.Config, opts ...Option) (sarama.SyncProducer, error) {
	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}
