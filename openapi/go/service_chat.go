package openapi

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    string
	SendQueue chan []byte
	RecvQueue chan []byte
}

type MessageQuery struct {
	client *Client
	data   []byte
}

type SimpleMessage struct {
	client *Client
	data   []byte
}

type ChatAPP struct {
	Clients      map[*Client]bool
	RecvAggQueue chan interface{}
	ConnQueue    chan *Client
	DisconnQueue chan *Client
}

func NewChatAPP() *ChatAPP {
	return &ChatAPP{
		Clients:      make(map[*Client]bool),
		RecvAggQueue: make(chan interface{}),
		ConnQueue:    make(chan *Client),
		DisconnQueue: make(chan *Client),
	}
}

func (c *ChatAPP) StartAPP() {
	go func() {
		c.RecvAggQueue = make(chan interface{})
		defer close(c.RecvAggQueue)

		for {
			select {
			case msg := <-c.RecvAggQueue:
				if v, ok := msg.(SimpleMessage); ok {
					resp := append([]byte(v.client.UserID+": "), v.data...)
					for client := range c.Clients {
						client.SendQueue <- resp
					}

					fmt.Println("recved value: ", v)
				}
			case client := <-c.ConnQueue:
				c.Clients[client] = true
				client.RecvQueue = make(chan []byte)
				client.SendQueue = make(chan []byte)
				go func() {
					for msg := range client.RecvQueue {
						c.RecvAggQueue <- SimpleMessage{client, msg}
					}
					fmt.Println("foward pump exit", client)
				}()
			case client := <-c.DisconnQueue:
				if _, ok := c.Clients[client]; ok {
					fmt.Println("closed", client)
					delete(c.Clients, client)
					close(client.SendQueue)
					close(client.RecvQueue)
				}
			}
		}
	}()
}

func (c *ChatAPP) receiver(client *Client, conn *websocket.Conn, userId string) {
	defer func() {
		conn.Close()
		c.DisconnQueue <- client
	}()

	for {
		fmt.Println("receiver waiting for message")

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("err = %+v\n", err)
			return
		}

		if messageType != websocket.TextMessage {
			fmt.Println("Unexpected message type:", messageType)
			return
		}

		client.RecvQueue <- p
	}
}

func (c *ChatAPP) sender(client *Client, conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		fmt.Println("sender select")
		select {
		case sendMsg, ok := <-client.SendQueue:
			if ok {
				conn.WriteMessage(websocket.TextMessage, sendMsg)
			} else {
				fmt.Println("sender exit")
				return
			}
		}
	}
}

func (c *ChatAPP) StartWebSocketHandler(conn *websocket.Conn, userId string) {
	client := Client{
		UserID:    userId,
		SendQueue: nil,
		RecvQueue: nil,
	}

	fmt.Println("starting service")
	c.ConnQueue <- &client
	go c.receiver(&client, conn, userId)
	go c.sender(&client, conn)
}
