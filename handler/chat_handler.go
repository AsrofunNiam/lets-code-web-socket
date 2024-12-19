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
	var msg domain.Message
	var roomID string
	var username string
	var isGroup bool

	if strings.HasPrefix(path, "group/") {
		// Group chat
		roomID = strings.TrimPrefix(path, "group/")
		isGroup = true
		msg.Group = true
	} else {
		// Private chat
		// Implementation room your take user name from start chat to target chat
		// Its sample /chat/{user1}_and_{user2} or /chat/{user2}_and_{user1}
		users := strings.Split(path, "_and_")
		if len(users) != 2 {
			http.Error(w, "Invalid room URL format. Use /chat/private/{user1}_and_{user2} or /chat/group/{groupID}", http.StatusBadRequest)
			return
		}
		roomID = chatHandler.ChatService.CreateRoomID(users[0], users[1])
		username = users[0] // Assume the first user is the sender
	}

	conn, err := upgraded.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	room := chatHandler.ChatService.ChatRepository.GetOrCreateRoom(roomID)
	room.Mutex.Lock()
	if !isGroup {
		room.Clients[conn] = username
	} else {
		// For group, associate the connection with a generic room ID or prompt for a user-specific identifier
		room.Clients[conn] = roomID
	}
	room.Mutex.Unlock()

	// Send history
	chatHandler.ChatService.SendHistory(roomID, conn)

	// Handle incoming messages and loop to read and broadcast messages
	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			// chatHandler.ChatService.ChatRepository.RemoveClientFromRoom(roomID, conn)
			go func() {
				chatHandler.ChatService.ChatRepository.RemoveClientFromRoom(roomID, conn)
			}()
			break
		}
		// Add timestamp and broadcast message
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		chatHandler.ChatService.BroadcastMessage(roomID, msg, conn)
	}
}
