package kafka

import (
	"chat-server/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg      *config.Config
	producer *kafka.Producer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{
		cfg: cfg,
	}

	var err error
	if k.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URL,
		"client.id":         cfg.Kafka.ClientID,
		"acks":              "all",
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}
