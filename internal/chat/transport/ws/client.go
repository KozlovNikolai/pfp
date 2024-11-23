// Package ws ...
package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Subscriber ...
type Subscriber struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	ChatID   string `json:"chat_id"`
	Username string `json:"username"`
}

// Message ...
type Message struct {
	Content  string `json:"content"`
	ChatID   string `json:"chat_id"`
	Username string `json:"username"`
}

func (c *Subscriber) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		m, ok := <-c.Message
		if !ok {
			return
		}
		err := c.Conn.WriteJSON(m)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func (c *Subscriber) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
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
			Content:  string(message),
			ChatID:   c.ChatID,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
