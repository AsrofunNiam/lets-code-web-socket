package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	// Jalankan server di port 8080
	fmt.Println("WebSocket server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

// Inisialisasi upgrader untuk meng-upgrade koneksi HTTP menjadi WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Mengizinkan semua origin, sesuaikan sesuai kebutuhan Anda
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade koneksi HTTP ke WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	// Loop untuk menerima dan mengirim pesan
	for {
		// Membaca pesan dari client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		// Mengirim kembali pesan ke client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
