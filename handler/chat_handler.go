package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AsrofunNiam/lets-code-web-socket/model/domain"
	"github.com/AsrofunNiam/lets-code-web-socket/service"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	ChatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{ChatService: chatService}
}

var upgraded = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle WebSocket connections
func (chatHandler *ChatHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/chat/")

	// implementation room your take user name from start chat to target chat
	// is sample /chat/{user1}_and_{user2} or /chat/{user2}_and_{user1}
	users := strings.Split(path, "_and_")

	// // Validate user accept max room or min room
	// if len(users) != 2 {
	// 	http.Error(w, "Invalid room URL format. Use /chat/{user1}_and_ ", http.StatusBadRequest)
	// 	return
	// }

	roomID := chatHandler.ChatService.CreateRoomID(users[0], users[1])

	conn, err := upgraded.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	room := chatHandler.ChatService.ChatRepository.GetOrCreateRoom(roomID)
	room.Mutex.Lock()
	room.Clients[conn] = users[0]
	fmt.Println("Client connected to room ascending:", roomID)
	fmt.Println("Clients started room:", users[0])
	room.Mutex.Unlock()

	// Handle incoming messages and loop to read and broadcast messages
	for {
		var msg domain.Message
		err := conn.ReadJSON(&msg)
		fmt.Println(" message read time:", msg.Timestamp)
		if err != nil {
			fmt.Println("Error reading message:", err)
			chatHandler.ChatService.ChatRepository.RemoveClientFromRoom(roomID, conn)
			break
		}

		// Broadcast the message
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		chatHandler.ChatService.BroadcastMessage(roomID, msg, conn)
	}
}
