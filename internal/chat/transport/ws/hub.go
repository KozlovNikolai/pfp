package ws

import "fmt"

// Chat ...
type Chat struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Clients map[string]*Subscriber `json:"client"`
}

// Hub ...
type Hub struct {
	Chats      map[string]*Chat
	Register   chan *Subscriber
	Unregister chan *Subscriber
	Broadcast  chan *Message
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		Chats:      make(map[string]*Chat),
		Register:   make(chan *Subscriber),
		Unregister: make(chan *Subscriber),
		Broadcast:  make(chan *Message, 5),
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// если комната, указанная в клиенте не существует, то выходим из селекта
			if _, ok := h.Chats[client.ChatID]; !ok {
				fmt.Printf("..Register..............Chat whith id: %s not exists\n", client.ChatID)
				break
			}
			// если комната существует, то работаем с ней:
			chat := h.Chats[client.ChatID]
			// если в комнате этот клиент не зарегистрировано, то регистрируем его:
			if _, ok := chat.Clients[client.ID]; !ok {
				fmt.Printf("..Register..............Registering client %s in chat %s\n", client.ID, client.ChatID)
				chat.Clients[client.ID] = client
			}
		case client := <-h.Unregister:
			// если комната, указанная в клиенте не существует, то выходим из селекта
			if _, ok := h.Chats[client.ChatID]; !ok {
				fmt.Printf("..Unregister..............Chat whith id: %s not exists\n", client.ChatID)
				break
			}
			// если в комнате, указанной в клиенте не зарегистрирован этот клиент, то выходим из селекта
			if _, ok := h.Chats[client.ChatID].Clients[client.ID]; !ok {
				fmt.Printf("..Unregister..............Client %s in chat %s not registering\n", client.ID, client.ChatID)
				break
			}
			// если в комнате, указанной в клиенте кто то еще остался, то отправляем им сообщение:
			if len(h.Chats[client.ChatID].Clients) != 0 {
				h.Broadcast <- &Message{
					ChatID:   client.ChatID,
					Content:  "user left the chat",
					Username: client.Username,
				}
				delete(h.Chats[client.ChatID].Clients, client.ID) // удаляем клиента из комнаты
				close(client.Message)
			}
		case message := <-h.Broadcast:
			// если комната, указанная в сообщении не существует, то выходим из селекта
			if _, ok := h.Chats[message.ChatID]; !ok {
				break
			}
			// идем по всем клиентам в комнате и отправляем им сообщение:
			for _, client := range h.Chats[message.ChatID].Clients {
				client.Message <- message
			}
		}
	}
}
