package network

import (
	"chat_socket_server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	Forward chan *message // 수신 메시지를 보관하는 채널
	// 들어오는 메시지를 다른 모든 클라이언트로 보내는 데 사용한다.

	Join  chan *Client // 접속에 대한 채널
	Leave chan *Client // 접속에 대한 채널

	Clients map[*Client]bool // 현재 방에 있는 모든 클라이언트를 의미

	service *service.Service
}

type Client struct {
	Socket *websocket.Conn // client의 웹 소켓
	Send   chan *message   // 전송되는 채널
	Room   *Room           // 클라이언트가 속해 있는 방
	Name   string          `json:"name"`
}

type message struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Room    string    `json:"room"`
	When    time.Time `json:"when"`
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		service: service,
	}
}

func (c *Client) Read() {
	// 클라이언트가 ReadMessage메소드를 통해서 소켓에서 읽을 수 있고,
	// 받은 메시지를 room타입에게 계속해서 전송을 한다.
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			return
		}

		msg.When = time.Now()
		msg.Name = c.Name
		c.Room.Forward <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true // client가 새로 들어 올떄
		case client := <-r.Leave:
			delete(r.Clients, client) // 나갈 떄에는 map값에서 client를 제거
			close(client.Send)        // 이후 client의 socker을 닫는다.
		case msg := <-r.Forward: // 만약 특정 메시지가 방에 들어오면

			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)

			for client := range r.Clients {
				client.Send <- msg // 모든 client에게 전달 해 준다.
			}
		}
	}
}

const (
	SocketBufferSize  = 1024
	messageBufferSize = 256
)

// 기본적으로 HTTP에 웹 소켓을 사용하려면, 이와 같이 업그레이드 해주어야 한다.
// -> 재사용 가능하기 떄문에 하나만 만들어도 된다.
var upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: messageBufferSize}

func (r *Room) ServeHTTP(c *gin.Context) {
	// 이후 요청이 이렇게 들어오게 된다면 Upgrade를 통해서 소켓을 가져 온다.

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

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

	// 문제가 없다면 client를 생성하여 방에 입장했다고 채널에 전송한다.
	client := &Client{
		Socket: Socket,
		Send:   make(chan *message, messageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- client

	// 또한 defer를 통해서 client가 끝날 떄를 대비하여 퇴장하는 작업을 연기한다.
	defer func() { r.Leave <- client }()

	// 이 후 고루틴을 통해서 write를 실행 시킨다.
	go client.Write()
	// 이 후 메인 루틴에서 read를 실행함으로써 해당 요청을 닫는것을 차단한다.
	// -> 연결을 활성화 시키는 것이다. 채널을 활용하여
	client.Read()
}
