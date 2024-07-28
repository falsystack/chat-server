package network

import "golang.org/x/net/websocket"

type Room struct {
	Forward chan *message    // 수신되는 메세지를 보관
	Join    chan *Client     // Socket 이 연결되는 경우에 작동
	Leave   chan *Client     // Socket 이 끊어지는 경우에 작동
	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}
