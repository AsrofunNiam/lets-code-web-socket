package domain

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	Clients map[*websocket.Conn]string
	Mutex   *sync.Mutex
	History []Message
	IsGroup bool
}
