package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// ChatRoom
type ChatRoom struct {
	clients map[*websocket.Conn]string
	mutex   *sync.Mutex
}

var rooms = make(map[string]*ChatRoom) // Map to all rooms
var globalMutex = &sync.Mutex{}        // Mutex for global access

// Unix RoomID
func createRoomID(user1, user2 string) string {
	users := []string{user1, user2}
	sort.Strings(users) // Urutkan users
	return "room_" + strings.Join(users, "_")
}

// Fungsi untuk menangani koneksi baru
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Ekstrak nama dari URL
	path := strings.TrimPrefix(r.URL.Path, "/chat/")
	users := strings.Split(path, "_and_")
	if len(users) != 2 {
		http.Error(w, "Invalid room URL format. Use /chat/{user1}_and_{user2}", http.StatusBadRequest)
		return
	}

	// Create room ID
	roomID := createRoomID(users[0], users[1])

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Move the connection to the room
	globalMutex.Lock()
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &ChatRoom{
			clients: make(map[*websocket.Conn]string),
			mutex:   &sync.Mutex{},
		}
	}
	globalMutex.Unlock()

	room := rooms[roomID]
	room.mutex.Lock()
	room.clients[conn] = "Anonymous" // Default username
	room.mutex.Unlock()

	// Loop to read and broadcast messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			room.mutex.Lock()
			delete(room.clients, conn)
			room.mutex.Unlock()
			break
		}

		// Send message to all clients
		broadcastMessage(msg, room)
	}
}

// Fungsi untuk mengirim pesan ke semua client
func broadcastMessage(msg Message, room *ChatRoom) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	for client := range room.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(room.clients, client)
		}
	}
}

func main() {
	// Dynamical routing
	http.HandleFunc("/chat/", handleConnections)

	fmt.Println("Chat server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
