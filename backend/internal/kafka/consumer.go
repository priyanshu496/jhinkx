package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"

	// Add your project imports here:
	"github.com/priyanshu496/jhinkx.git/internal/config"
	"github.com/priyanshu496/jhinkx.git/internal/db"
	"github.com/priyanshu496/jhinkx.git/internal/models"
)

var matchQueues = make(map[int][]string)
var queueMutex = &sync.Mutex{}

func InitConsumer(cfg *config.AppConfig) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    "matchmaking-requests",
		GroupID:  "matchmaking-workers",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("Matchmaking Consumer is online! Waiting for tickets...")

	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[KAFKA CONSUMER ERROR] Failed to read message: %v\n", err)
				break
			}

			var req MatchRequest
			if err := json.Unmarshal(m.Value, &req); err != nil {
				log.Printf("Failed to parse ticket data: %v\n", err)
				continue
			}

			log.Printf("🎟️ User %s entered the waiting room for size %d\n", req.UserID, req.PreferredGroupSize)
			processMatchmaking(req)
		}
	}()
}

func processMatchmaking(req MatchRequest) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	size := req.PreferredGroupSize

	for _, existingUserID := range matchQueues[size] {
		if existingUserID == req.UserID {
			log.Printf("⚠️ User %s is already in line! Ignoring duplicate ticket.\n", req.UserID)
			return
		}
	}


	matchQueues[size] = append(matchQueues[size], req.UserID)

	if len(matchQueues[size]) >= size {
		matchedUsers := matchQueues[size][:size]
		matchQueues[size] = matchQueues[size][size:]

		log.Printf("\n🎉 SUCCESS! GROUP FORMED! 🎉\nSize: %d\nUsers: %v\n", size, matchedUsers)

		// --- THE FINAL PIECE: POSTGRESQL INSERTION ---

		// 1. Create the new Space using your exact model schema
		newSpace := models.Space{
			TargetSize: size,
			Status:     "active",
		}

		if err := db.DB.Create(&newSpace).Error; err != nil {
			log.Printf("Failed to create space in DB: %v\n", err)
			return
		}

		// 2. Loop through our matched users and securely link them to the new Space
		for _, userID := range matchedUsers {
			member := models.SpaceMember{
				SpaceID: newSpace.ID,
				UserID:  userID,
			}
			if err := db.DB.Create(&member).Error; err != nil {
				log.Printf("Failed to add user %s to space: %v\n", userID, err)
			}
		}

		log.Println("✅ Successfully saved new Matchmaking Space to PostgreSQL!")
	} else {
		log.Printf("⏳ Waiting for %d more users to form a group of %d...\n", size-len(matchQueues[size]), size)
	}
}