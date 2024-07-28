package network

import (
	. "chat-server/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// http -> websocket 으로 업그레이드
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	Forward chan *message    // 수신되는 메세지를 보관
	Join    chan *client     // Socket 이 연결되는 경우에 작동
	Leave   chan *client     // Socket 이 끊어지는 경우에 작동
	Clients map[*client]bool // 현재 방에 있는 client 정보를 저장
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		Clients: make(map[*client]bool),
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &client{
		Socket: socket,
		Send:   make(chan *message, MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client

	defer func() {
		r.Leave <- client
	}()
}
