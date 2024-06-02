package network

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	SocketBufferSize  = 1024
	messageBufferSize = 256
)

// Upgrader HTTP -> Websocket　にアップグレードするときに使用する
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: messageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Room chat room
type Room struct {
	Forward chan *message    // 受信されるメッセージを保存、入ってくるメッセージを他のクライアントに転送
	Join    chan *Client     // Socket が繋がる場合動く
	Leave   chan *Client     // Socket がきれる場合動く
	Clients map[*Client]bool // 現在の Room にある Client の情報を保存
}

type message struct {
	Name    string
	Message string
	When    time.Time
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

func (c *Client) Read() {
	// Client　が入ってくるメッセージを読みとるメソッド
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
		}
		msg.When = time.Now()
		msg.Name = c.Name

		// 受け取ったメッセージを roomタイプに連続的に転送する
		c.Room.Forward <- msg
	}
}

func (c *Client) Write() {
	// Client　がメッセージを転送するメソッド
	defer c.Socket.Close()
	for msg := range c.Send {
		log.Println("WRITE : ", msg, "Client", c.Name)
		log.Println()
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) Run() {
	// Room にある全てのchanの値を受けとる役割
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			//r.Clients[client] = false
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			for client := range r.Clients {
				// 全ての client にメッセージを送る
				client.Send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	Socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("---- serveHTTP:", err)
		return
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("auth cookie is failed", err)
		return
	}

	// 問題がなければ　Client　を生成して Room に入場したお知らせを chan に送る
	client := &Client{
		Send:   make(chan *message, messageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
		Socket: Socket,
	}

	r.Join <- client
	// defer を利用して Client が終了する時退場させる
	defer func() { r.Leave <- client }()
	// go routine を利用して Writeを実行
	go client.Write()

	// そのあと main routine で Read を実行することで該当するリクエストが閉じられるのを防ぐ
	// -> channel を利用して連結を活性化させることである
	client.Read()
}
