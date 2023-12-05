package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	//list of all the clients
	clients map[*client]bool

	//channel join for adding new clients to the room
	join chan *client

	//channel forward for forwarding a new message to all the users
	forward chan []byte

	//channel leave, to remove a client from the room
	leave chan *client
}

// create a new chat room
func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		join:    make(chan *client),
		forward: make(chan []byte),
		leave:   make(chan *client),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)

		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		return
	}

	client := &client{
		socket:  *socket,
		receive: make(chan []byte),
		Room:    r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
