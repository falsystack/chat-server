package network

import "C"
import (
	"chat-server/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// Upgrader HTTP -> Websocket　にアップグレードするときに使用する
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Room chat room
type Room struct {
	Forward chan *message    // 受診されるメッセージを保存、入ってくるメッセージを他のクライアントに転送
	Join    chan *Client     // Socket が繋がる場合動く
	Leave   chan *Client     // Socket がきれる場合動く
	Clients map[*Client]bool // 現在の Room にある Client の情報を保存
}

type message struct {
	Name    string
	Message string
	TIme    int64
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) RunInit() {
	// Room にある全てのchanの値を受けとる役割
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// TODO: panic 使用はよくないので修正
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		// TODO: panic 使用はよくないので修正
		panic(err)
	}

	client := &Client{
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
		Socket: socket,
	}

	r.Join <- client
	defer func() { r.Leave <- client }()

	//
}
