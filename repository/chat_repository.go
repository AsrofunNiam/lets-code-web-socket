package repository

import (
	"fmt"
	"sync"

	"github.com/AsrofunNiam/lets-code-web-socket/model/domain"
	"github.com/gorilla/websocket"
)

type ChatRepository struct {
	Rooms map[string]*domain.ChatRoom
	Mutex *sync.Mutex
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		Rooms: make(map[string]*domain.ChatRoom),
		Mutex: &sync.Mutex{},
	}
}

// Get or create room
func (repo *ChatRepository) GetOrCreateRoom(roomID string) *domain.ChatRoom {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	test := &repo.Rooms

	fmt.Println("Get or create room:", roomID)
	fmt.Println("Alamat memori num repo:", &repo.Rooms)
	fmt.Println("Alamat memori num var:", test)
	fmt.Println("Length var:", len(*test))
	fmt.Println("Length repo:", len(repo.Rooms))

	if _, exists := repo.Rooms[roomID]; !exists {
		repo.Rooms[roomID] = &domain.ChatRoom{
			Clients: make(map[*websocket.Conn]string),
			Mutex:   &sync.Mutex{},
		}
	}
	return repo.Rooms[roomID]
}

// Remove client from room
func (repo *ChatRepository) RemoveClientFromRoom(roomID string, conn *websocket.Conn) {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	if room, exists := repo.Rooms[roomID]; exists {
		room.Mutex.Lock()
		delete(room.Clients, conn)
		room.Mutex.Unlock()
	}
}
