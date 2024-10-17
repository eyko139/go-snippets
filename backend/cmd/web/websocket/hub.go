package websocket

import "fmt"

type Hub struct {
	Clients map[*Client]bool
	// incoming messages
	broadcast            chan []byte
	broadcastToRecipient chan *wsMessage
	// register requests from clients
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:            make(chan []byte),
		broadcastToRecipient: make(chan *wsMessage),
		Clients:              make(map[*Client]bool),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case aimedMessage := <-h.broadcastToRecipient:
			for client := range h.Clients {
				if client.Id.String() == aimedMessage.Recipient {
                    fmt.Println(aimedMessage.Recipient)
					select {
					case client.send <- []byte(aimedMessage.Message):
					default:
						close(client.send)
						delete(h.Clients, client)
					}
				}
			}

		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
