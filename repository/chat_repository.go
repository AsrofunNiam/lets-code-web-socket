package repository

import (
	"fmt"
	"sync"

	"github.com/AsrofunNiam/lets-code-web-socket/model/domain"
	"github.com/gorilla/websocket"
)

type ChatRepository struct {
	ChatRoom map[string]*domain.ChatRoom
	Mutex    *sync.Mutex
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		ChatRoom: make(map[string]*domain.ChatRoom),
		Mutex:    &sync.Mutex{},
	}
}

// Get or create room
func (chatRepository *ChatRepository) GetOrCreateRoom(roomID string) *domain.ChatRoom {
	chatRepository.Mutex.Lock()
	defer chatRepository.Mutex.Unlock()

	// test := &chatRepository.ChatRoom
	// fmt.Println("Get or create room:", roomID)
	// fmt.Println("Alamat memori num repo:", &chatRepository.ChatRoom)
	// fmt.Println("Alamat memori num var:", test)
	// fmt.Println("Length var:", len(*test))
	// fmt.Println("Length repo:", len(chatRepository.ChatRoom))

	if _, exists := chatRepository.ChatRoom[roomID]; !exists {
		chatRepository.ChatRoom[roomID] = &domain.ChatRoom{
			Clients: make(map[*websocket.Conn]string),
			Mutex:   &sync.Mutex{},
		}
	}
	return chatRepository.ChatRoom[roomID]
}

// Remove client from room
func (chatRepository *ChatRepository) RemoveClientFromRoom(roomID string, conn *websocket.Conn) {
	chatRepository.Mutex.Lock()
	defer chatRepository.Mutex.Unlock()

	// fmt.Println("Memory num repo start:", &chatRepository.ChatRoom)
	// fmt.Println("Length repo: start", len(chatRepository.ChatRoom))

	if room, exists := chatRepository.ChatRoom[roomID]; exists {
		room.Mutex.Lock()
		delete(room.Clients, conn)

		// If implement delete room
		if len(room.Clients) == 0 {

			// // Implement Backup room history to redis or database
			// repo.BackupRoomHistoryToDB(roomID, room.History)

			delete(chatRepository.ChatRoom, roomID)
			fmt.Printf("Room %s deleted\n because no more clients\n", roomID)
		}
		room.Mutex.Unlock()
	}

	// fmt.Println("Memory num repo num repo end:", &chatRepository.ChatRoom)
	// fmt.Println("Length repo: end", len(chatRepository.ChatRoom))
}
