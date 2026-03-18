package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/priyanshu496/jhinkx.git/internal/config"
)

var MatchWriter *kafka.Writer

// MatchRequest is the structure of the data we will drop into the Kafka queue
type MatchRequest struct {
	UserID             string `json:"user_id"`
	PreferredGroupSize int    `json:"preferred_group_size"`
}

// InitProducer connects to our local Docker Kafka
func InitProducer(cfg *config.AppConfig) {
	// Using the modern struct initialization for kafka-go
	MatchWriter = &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaBroker),
		Topic:                  "matchmaking-requests",
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true, // THIS IS THE MAGIC KEY!
	}

	log.Println("Successfully connected to Local Kafka Producer!")
}

// PublishMatchRequest converts our struct to JSON and sends it to Kafka
func PublishMatchRequest(userID string, groupSize int) error {
	req := MatchRequest{
		UserID:             userID,
		PreferredGroupSize: groupSize,
	}

	// Convert the struct to a JSON string
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Send it to the Kafka topic
	err = MatchWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(userID), // Using UserID as the key keeps things organized
			Value: bytes,
		},
	)
	return err
}