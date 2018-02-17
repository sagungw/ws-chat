package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Channel struct {
	Clients  map[*websocket.Conn]bool
	Messages chan interface{}
}

func NewChannel() *Channel {
	return &Channel{
		Clients:  make(map[*websocket.Conn]bool),
		Messages: make(chan interface{}),
	}
}

func (c *Channel) Listen() {
	go c.listen()
}

func (c *Channel) RegisterClient(client *websocket.Conn) {
	c.Clients[client] = true
}

func (c *Channel) UnregisterClient(client *websocket.Conn) {
	delete(c.Clients, client)
}

func (c *Channel) listen() {
	for message := range c.Messages {
		for client := range c.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				c.UnregisterClient(client)
			}
		}
	}
}

type UserWSChannel struct {
	Clients map[*websocket.Conn]bool
	Users   chan *User
}

func NewUserWSChannel() *UserWSChannel {
	return &UserWSChannel{
		Clients: make(map[*websocket.Conn]bool),
		Users:   make(chan *User),
	}
}

func (c *UserWSChannel) Listen() {
	go c.listen()
}

func (c *UserWSChannel) RegisterClient(client *websocket.Conn) {
	c.Clients[client] = true
}

func (c *UserWSChannel) UnregisterClient(client *websocket.Conn) {
	delete(c.Clients, client)
}

func (c *UserWSChannel) listen() {
	for user := range c.Users {
		InsertUserIfNotExist(user)

		for client := range c.Clients {
			err := client.WriteJSON(Users)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				c.UnregisterClient(client)
			}
		}
	}
}
