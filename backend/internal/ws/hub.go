package ws

// BroadcastMessage wraps the JSON data with the specific SpaceID it belongs to
type BroadcastMessage struct {
	SpaceID string
	Data    []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients (the switchboard holding all the open cables)
	clients map[*Client]bool

	// Inbound messages from the clients waiting to be broadcast
	broadcast chan *BroadcastMessage // Changed to our new struct

	// Register requests from newly connecting clients
	register chan *Client

	// Unregister requests from clients who disconnect or close their browser
	unregister chan *Client
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the Hub in the background to continuously listen for traffic
func (h *Hub) Run() {
	for {
		select {
		// 1. A NEW USER CONNECTED!
		case client := <-h.register:
			h.clients[client] = true

		// 2. A USER DISCONNECTED!
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // Remove them from the active list
				close(client.send)        // Cut the cable
			}

		// 3. A NEW MESSAGE ARRIVED!
		case message := <-h.broadcast:
			// ROOM ISOLATION: Only send to users in the exact same Space!
			for client := range h.clients {
				if client.spaceID == message.SpaceID {
					select {
					case client.send <- message.Data:
						// Message sent successfully!
					default:
						// If the send buffer is full/blocked, assume the connection is dead
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}