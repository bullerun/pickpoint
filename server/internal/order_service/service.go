package order_service

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"OzonHW1/server/internal/storage/postgres"
	"github.com/IBM/sarama"
)

const (
	acceptReturn = "accept-return"
	acceptOrder  = "accept-order"
	returnOrder  = "return-order"
	updateIssued = "issue-order"
)

type Implementation struct {
	storage postgres.OrderService
	order_service.UnimplementedOrderServiceServer
	producer sarama.SyncProducer
	topic    string
}

func NewImplementation(storage postgres.OrderService, prod sarama.SyncProducer, topic string) *Implementation {
	return &Implementation{storage: storage, producer: prod, topic: topic}
}
