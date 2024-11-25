// Package ws ...
package ws

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Subscriber ...
type Subscriber struct {
	Conn     *websocket.Conn
	Message  chan *MessageOne
	ID       int       `json:"id"`
	Pubsub   uuid.UUID `json:"pubsub"`
	Username string    `json:"username"`
}

// Message ...
type Message struct {
	Content     string `json:"content"`
	ChatID      int    `json:"chat_id"`
	Sender      int    `json:"sender_id"`
	ChatMembers []int  `json:"chat_members"`
}
type MessageOne struct {
	Content string `json:"content"`
	ChatID  int    `json:"chat_id"`
	Sender  int    `json:"sender_id"`
}

func toMsgOne(msg Message) *MessageOne {
	return &MessageOne{
		Content: msg.Content,
		ChatID:  msg.ChatID,
		Sender:  msg.Sender,
	}
}

func (ss *Subscriber) writeMessage() {
	defer func() {
		ss.Conn.Close()
	}()
	for {
		m, ok := <-ss.Message
		if !ok {
			return
		}
		err := ss.Conn.WriteJSON(m)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func (ss *Subscriber) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- ss
		ss.Conn.Close()
	}()
	for {
		_, message, err := ss.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := &Message{
			Content: string(message),
			Sender:  ss.ID,
		}
		hub.Broadcast <- msg
	}
}
