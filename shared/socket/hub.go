package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// SSEClient represents a connected SSE client
type SSEClient struct {
	ID     string
	UserID string
	RoleID string
	Chan   chan SSEMessage
}

// SSEMessage represents a message to be sent via SSE
type SSEMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// SocketHub manages SSE connections (kept name for backward compatibility)
type SocketHub struct {
	clients    map[string]*SSEClient
	mu         sync.RWMutex
	register   chan *SSEClient
	unregister chan *SSEClient
}

var (
	instance *SocketHub
	once     sync.Once
)

// GetInstance returns the singleton instance of SocketHub
func GetInstance() *SocketHub {
	once.Do(func() {
		instance = &SocketHub{
			clients:    make(map[string]*SSEClient),
			register:   make(chan *SSEClient),
			unregister: make(chan *SSEClient),
		}
		go instance.run()
	})
	return instance
}

func (h *SocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf(">>> SSE Client Connected: %s (User: %s, Role: %s)", client.ID, client.UserID, client.RoleID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				close(client.Chan)
				delete(h.clients, client.ID)
				log.Printf("<<< SSE Client Disconnected: %s", client.ID)
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastToRole sends a message to all users in a specific role room
func (h *SocketHub) BroadcastToRole(roleID string, event string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SSEMessage{Event: event, Data: data}
	count := 0
	for _, client := range h.clients {
		if client.RoleID == roleID {
			select {
			case client.Chan <- msg:
				count++
			default:
				// Channel full, skip
			}
		}
	}
	log.Printf("Broadcast to role %s: event=%s, delivered to %d clients", roleID, event, count)
}

// SendToUser sends a message to a specific user's personal channel
func (h *SocketHub) SendToUser(userID string, event string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SSEMessage{Event: event, Data: data}
	count := 0
	for _, client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Chan <- msg:
				count++
			default:
			}
		}
	}
	log.Printf("Send to user %s: event=%s, delivered to %d clients", userID, event, count)
}

// ServeHTTP handles the SSE endpoint
func (h *SocketHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get user info from query params
	userID := r.URL.Query().Get("userId")
	roleID := r.URL.Query().Get("roleId")

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Create client
	client := &SSEClient{
		ID:     fmt.Sprintf("%s-%d", userID, r.Context().Value(http.LocalAddrContextKey)),
		UserID: userID,
		RoleID: roleID,
		Chan:   make(chan SSEMessage, 64),
	}

	// Use remote addr as unique ID
	client.ID = r.RemoteAddr

	h.register <- client

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Listen for client disconnect
	ctx := r.Context()

	// Add ticker for heartbeat (keepalive)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.unregister <- client
			return
		case <-ticker.C:
			// Send comment to keep connection alive
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		case msg, ok := <-client.Chan:
			if !ok {
				return
			}
			data, err := json.Marshal(msg.Data)
			if err != nil {
				data = []byte("{}")
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Event, string(data))
			flusher.Flush()
		}
	}
}
