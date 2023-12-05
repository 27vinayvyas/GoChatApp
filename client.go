package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	//socket is the connection for web socket
	socket websocket.Conn

	//room in which the user is chatting in
	Room *room

	//to receive messages from other client
	receive chan []byte
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.Room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
