package main

import (
	"fmt"
	"net/http"

	"github.com/AsrofunNiam/lets-code-web-socket/handler"
	"github.com/AsrofunNiam/lets-code-web-socket/repository"
	"github.com/AsrofunNiam/lets-code-web-socket/service"
)

func main() {
	// Initialize repository
	chatRepo := repository.NewChatRepository()

	// Initialize service
	chatService := service.NewChatService(chatRepo)

	// Initialize handler
	chatHandler := handler.NewChatHandler(chatService)

	// Set routes
	http.HandleFunc("/chat/", chatHandler.HandleConnections)

	fmt.Println("Chat server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
