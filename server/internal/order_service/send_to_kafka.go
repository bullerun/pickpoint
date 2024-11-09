package order_service

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func (s *Implementation) sendToKafka(event any, commandName string, orderID int64) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal event message")
	}
	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(orderID, 10)),
		Value: sarama.ByteEncoder(eventData),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("command"),
				Value: []byte(commandName),
			},
		},
		Timestamp: time.Now(),
	}
	_, _, err = s.producer.SendMessage(msg)
	if err != nil {
		return status.Error(codes.Internal, "failed to send Kafka message")
	}
	return nil
}
