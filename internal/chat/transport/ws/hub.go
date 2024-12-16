package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	Node         map[uuid.UUID]*Subscriber
	Register     chan *Subscriber
	Unregister   chan *Subscriber
	Broadcast    chan *Message
	stateService IStateService
}

// NewHub ...
func NewHub(stateService IStateService) *Hub {
	return &Hub{
		// Node:         make(map[int]*Subscriber),
		Node:         make(map[uuid.UUID]*Subscriber),
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
		case subscriber := <-h.Register:
			node := h.Node
			if _, ok := node[subscriber.Pubsub]; !ok {
				fmt.Printf("..Register..............Registering client %d\n", subscriber.ID)
				node[subscriber.Pubsub] = subscriber
				fmt.Println("############################################################################################################")
			}
			h.printConnections()
		case subscriber := <-h.Unregister:
			if _, ok := h.Node[subscriber.Pubsub]; !ok {
				log.Printf("..Unregister..............Client %d not registering\n", subscriber.ID)
				//break
			} else {
				log.Printf("..Unregister..............Client %d left the socket\n", subscriber.ID)
				delete(h.Node, subscriber.Pubsub)
				close(subscriber.Message)
				state, ok := h.stateService.DeleteConnFromState(context.Background(), subscriber.ID, subscriber.Pubsub)
				log.Printf("State after Unsubscribe: %+v, is OK: %v", state, ok)
				fmt.Println("############################################################################################################")
			}
			h.printConnections()
		case message := <-h.Broadcast:
			for _, userID := range message.ChatMembers {
				state, _ := h.stateService.GetState(context.Background(), userID)
				for _, connect := range state.Connects {
					ss, ok := h.Node[connect.Pubsub]
					if ok {
						ss.Message <- toMsgOne(*message)
					}
				}
			}
		}
	}
}

// // Run ...
// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case subscriber := <-h.Register:
// 			node := h.Node
// 			if _, ok := node[subscriber.Pubsub]; !ok {
// 				fmt.Printf("..Register..............Registering client %d\n", subscriber.ID)
// 				node[subscriber.Pubsub] = subscriber
// 				fmt.Println("############################################################################################################")
// 			}
// 			h.printConnections()
// 		case subscriber := <-h.Unregister:
// 			if _, ok := h.Node[subscriber.Pubsub]; !ok {
// 				log.Printf("..Unregister..............Client %d not registering\n", subscriber.ID)
// 				//break
// 			} else {
// 				log.Printf("..Unregister..............Client %d left the socket\n", subscriber.ID)
// 				delete(h.Node, subscriber.Pubsub)
// 				close(subscriber.Message)
// 				state, ok := h.stateService.DeleteConnFromState(context.Background(), subscriber.ID, subscriber.Pubsub)
// 				log.Printf("State after Unsubscribe: %+v, is OK: %v", state, ok)
// 				fmt.Println("############################################################################################################")
// 			}
// 			h.printConnections()
// 		case message := <-h.Broadcast:
// 			for _, userID := range message.ChatMembers {
// 				state, _ := h.stateService.GetState(context.Background(), userID)
// 				for _, connect := range state.Connects {
// 					ss, ok := h.Node[connect.Pubsub]
// 					if ok {
// 						ss.Message <- toMsgOne(*message)
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func (h *Hub) printConnections() {
	for _, v := range h.Node {
		fmt.Printf("user id: %d, pubsub: %v\n", v.ID, v.Pubsub)
	}
}
