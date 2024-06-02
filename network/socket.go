package network

import "C"
import (
	"chat-server/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// Upgrader HTTP -> Websocket　にアップグレードするときに使用する
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Room chat room
type Room struct {
	Forward chan *message    // 受信されるメッセージを保存、入ってくるメッセージを他のクライアントに転送
	Join    chan *client     // Socket が繋がる場合動く
	Leave   chan *client     // Socket がきれる場合動く
	Clients map[*client]bool // 現在の Room にある client の情報を保存
}

type message struct {
	Name    string
	Message string
	TIme    int64
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

func (c *client) Read() {
	// client　が入ってくるメッセージを読みとるメソッド
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			log.Println("READ : ", msg, "client", c.Name)
			log.Println()
			msg.TIme = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c client) Write() {
	// client　がメッセージを転送するメソッド
	defer c.Socket.Close()
	for msg := range c.Send {
		log.Println("WRITE : ", msg, "client", c.Name)
		log.Println()
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
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

	client := &client{
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
		Socket: socket,
	}

	r.Join <- client
	defer func() { r.Leave <- client }()

	go client.Write()
	client.Read()
}
