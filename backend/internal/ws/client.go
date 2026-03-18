package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/priyanshu496/jhinkx.git/internal/db"
	"github.com/priyanshu496/jhinkx.git/internal/models"
	"github.com/priyanshu496/jhinkx.git/internal/redis"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// WARNING: In production, check the origin to prevent CORS attacks!
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client represents a single user's live connection
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	spaceID string
	userID  string
}

// readPump pumps messages FROM the websocket connection TO the hub.
// readPump pumps messages FROM the websocket connection TO the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		redis.Client.SRem(redis.Ctx, "space:"+c.spaceID+":online", c.userID)
		log.Printf("User %s disconnected from space %s", c.userID, c.spaceID)
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 1. We expect the frontend to send us a JSON message like {"content": "Hello!"}
		var incoming struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(messageBytes, &incoming); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// 2. Save the message to PostgreSQL permanently!
		newMsg := models.Message{
			SpaceID: c.spaceID,
			UserID:  c.userID,
			Content: incoming.Content,
		}

		if err := db.DB.Create(&newMsg).Error; err != nil {
			log.Printf("Failed to save message to DB: %v", err)
			continue
		}

		// 3. Convert the official database record (with ID and CreatedAt) back to JSON
		broadcastBytes, _ := json.Marshal(newMsg)

		// 4. Send it to the Hub to be broadcasted to everyone in the room!
		c.hub.broadcast <- broadcastBytes
	}
}

// writePump pumps messages FROM the hub TO the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Send a ping to the client every 54 seconds to keep the connection alive
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS handles the initial HTTP request and upgrades it to a WebSocket
func ServeWS(hub *Hub, c *gin.Context) {
	spaceID := c.Param("id")
	// Note: You would normally get UserID from your JWT middleware here.
	// For testing, we'll grab it from a query param like ?userId=123
	userID := c.Query("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	// 1. Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}

	// 2. Create the Client
	client := &Client{
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		spaceID: spaceID,
		userID:  userID,
	}

	// 3. Plug them into the Hub
	client.hub.register <- client

	// 4. Mark them as ONLINE in Redis! (Using a Redis Set)
	redis.Client.SAdd(redis.Ctx, "space:"+spaceID+":online", userID)
	log.Printf("User %s connected to space %s", userID, spaceID)

	// 5. Start the pumps in background Go routines
	go client.writePump()
	go client.readPump()
}
