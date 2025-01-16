package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	//"golang.org/x/text/message"
)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, err := json.Marshal(&Message{Content: "socket has been connected 🥷"})
			if err != nil {
				fmt.Println("shuckles something went wrong")
				panic(err)

			}
			manager.send(jsonMessage, conn)

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "socket has been disconnected, bye bye"})
				manager.send(jsonMessage, conn)
			}
		// this case is for sending message for every client to see
		//
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}

		}
	}
}

// send message to everyone except messagee
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	for {
		_, message, err := c.socket.ReadMessage()
		// in case we can't read the message, we will unregister client
		// and close his connection
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}

		// convering the bytes into message json fomrat, and omitting recipient
		jsonMsg, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		// we add that json message to the messages that manager has to broadcast.
		manager.broadcast <- jsonMsg
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

}
