package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Room struct {
	Forward chan *message    // 수신되는 메세지를 보관
	Join    chan *Client     // Socket 이 연결되는 경우에 작동
	Leave   chan *Client     // Socket 이 끊어지는 경우에 작동
	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

type Client struct {
	Socket *websocket.Conn
	Send   chan *message
	Room   *Room
	Name   string
}

type message struct {
	Name    string
	Message string
	When    time.Time
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (c Client) Read() {
	// 클라이언트가 들어오는 메시지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg *message
		if err := c.Socket.ReadJSON(&msg); err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			msg.When = time.Now()
			msg.Name = c.Name
			c.Room.Forward <- msg
		}
	}
}

func (c Client) Write() {
	// 클라이언트가 메시지를 전송하는 함수
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}

}

func (r *Room) Run() {
	// Room 에 있는 모든 채널값들을 받는 역할을 한다
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

// http -> websocket 으로 업그레이드
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Room) ServeHTTP(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("----- serveHTTP", err)
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("auth cookie is failed", err)
		return
	}

	client := &Client{
		Socket: socket,
		Send:   make(chan *message, MessageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- client

	defer func() {
		r.Leave <- client
	}()

	go client.Write()
	client.Read()
}
