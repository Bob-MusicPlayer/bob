package util

// Hub maintains the set of active Clients and broadcasts messages to the
// Clients.
type Hub struct {
	// Registered Clients.
	Clients map[*Channel]bool

	// Register requests from the Clients.
	register chan *Channel

	// Unregister requests from Clients.
	unregister chan *Channel

	OnRegister chan *Channel

	OnUnregister chan *Channel
}

func NewHub() *Hub {
	return &Hub{
		register:     make(chan *Channel),
		unregister:   make(chan *Channel),
		OnRegister:   make(chan *Channel),
		OnUnregister: make(chan *Channel),
		Clients:      make(map[*Channel]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
			h.OnRegister <- client
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				h.OnUnregister <- nil
			}
		}
	}
}

func (h *Hub) Broadcast(packet *Packet) {
	for client := range h.Clients {
		client.Send <- packet
	}
}
