package ws

import "fmt"

type Hub struct {
	Node         map[int]*Subscriber
	Register     chan *Subscriber
	Unregister   chan *Subscriber
	Broadcast    chan *Message
	stateService IStateService
}

// NewHub ...
func NewHub(stateService IStateService) *Hub {
	return &Hub{
		Node:         make(map[int]*Subscriber),
		Register:     make(chan *Subscriber),
		Unregister:   make(chan *Subscriber),
		Broadcast:    make(chan *Message, 5),
		stateService: stateService,
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			node := h.Node
			if _, ok := node[client.ID]; !ok {
				fmt.Printf("..Register..............Registering client %d\n", client.ID)
				node[client.ID] = client
			}
		case client := <-h.Unregister:
			if _, ok := h.Node[client.ID]; !ok {
				fmt.Printf("..Unregister..............Client %d not registering\n", client.ID)
				break
			}
		case message := <-h.Broadcast:
			for _, client := range h.Node {
				client.Message <- message
			}
		}
	}
}
