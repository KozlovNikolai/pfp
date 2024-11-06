package ws

import "fmt"

// Room ...
type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"client"`
}

// Hub ...
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// если комната, указанная в клиенте не существует, то выходим из селекта
			if _, ok := h.Rooms[client.RoomID]; !ok {
				fmt.Printf("..Register..............Room whith id: %s not exists\n", client.RoomID)
				break
			}
			// если комната существует, то работаем с ней:
			room := h.Rooms[client.RoomID]
			// если в комнате этот клиент не зарегистрировано, то регистрируем его:
			if _, ok := room.Clients[client.ID]; !ok {
				fmt.Printf("..Register..............Registering client %s in room %s\n", client.ID, client.RoomID)
				room.Clients[client.ID] = client
			}
		case client := <-h.Unregister:
			// если комната, указанная в клиенте не существует, то выходим из селекта
			if _, ok := h.Rooms[client.RoomID]; !ok {
				fmt.Printf("..Unregister..............Room whith id: %s not exists\n", client.RoomID)
				break
			}
			// если в комнате, указанной в клиенте не зарегистрирован этот клиент, то выходим из селекта
			if _, ok := h.Rooms[client.RoomID].Clients[client.ID]; !ok {
				fmt.Printf("..Unregister..............Client %s in room %s not registering\n", client.ID, client.RoomID)
				break
			}
			// если в комнате, указанной в клиенте кто то еще остался, то отправляем им сообщение:
			if len(h.Rooms[client.RoomID].Clients) != 0 {
				h.Broadcast <- &Message{
					RoomID:   client.RoomID,
					Content:  "user left the chat",
					Username: client.Username,
				}
				delete(h.Rooms[client.RoomID].Clients, client.ID) // удаляем клиента из комнаты
				close(client.Message)
			}
		case message := <-h.Broadcast:
			// если комната, указанная в сообщении не существует, то выходим из селекта
			if _, ok := h.Rooms[message.RoomID]; !ok {
				break
			}
			// идем по всем клиентам в комнате и отправляем им сообщение:
			for _, client := range h.Rooms[message.RoomID].Clients {
				client.Message <- message
			}
		}
	}
}
