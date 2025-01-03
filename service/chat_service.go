package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AsrofunNiam/lets-code-web-socket/model/domain"
	"github.com/AsrofunNiam/lets-code-web-socket/repository"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	ChatRepository *repository.ChatRepository
}

func NewChatService(
	chatRepository *repository.ChatRepository,
) *ChatService {
	return &ChatService{ChatRepository: chatRepository}
}

// Generate unique RoomID and sort room name ()
func (chatService *ChatService) CreateRoomID(user1, user2 string) string {
	users := []string{user1, user2}
	sort.Strings(users)
	return "room_" + strings.Join(users, "_")
}

// Broadcast message to room
func (chatService *ChatService) BroadcastMessage(roomID string, msg domain.Message, sender *websocket.Conn) {
	room := chatService.ChatRepository.GetOrCreateRoom(roomID)

	// Message history
	room.Mutex.Lock()
	room.History = append(room.History, msg)
	defer room.Mutex.Unlock()

	fmt.Println("message broadcast time:", msg.Timestamp)
	fmt.Println("obj broadcast room:", room)
	fmt.Println("obj broadcast room", room.History)

	for client := range room.Clients {

		// // ignore sender
		// if client == sender {
		// 	continue
		// }

		err := client.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(room.Clients, client)
		}
	}
}
