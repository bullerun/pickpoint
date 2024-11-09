package consumer_group

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	ready chan bool
}

func NewConsumerGroupHandler(
// map[topic]TopicHandler
) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор, пока сессия не завершилась
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := convertMsg(message)
			data, _ := json.Marshal(msg)
			log.Printf("Message claimed: %s", data)

			session.MarkMessage(message, "")

			session.Commit()

		case <-session.Context().Done():
			return nil
		}
	}
}

type Msg struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func convertMsg(in *sarama.ConsumerMessage) Msg {
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       string(in.Key),
		Payload:   string(in.Value),
	}
}