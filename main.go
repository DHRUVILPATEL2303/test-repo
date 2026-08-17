package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// We need an upgrader to upgrade the HTTP connection to a WebSocket
var upgrader = websocket.Upgrader{
	// Allow all origins for testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected!")

	// Read the test environment variables
	testEnv1 := os.Getenv("TEST_ENV_1")
	testEnv2 := os.Getenv("TEST_ENV_2")

	// Send a welcome message containing the environment variables
	welcomeMsg := fmt.Sprintf("Welcome to the Go Echo WebSocket Server!\nTEST_ENV_1: %s\nTEST_ENV_2: %s", testEnv1, testEnv2)
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// Echo loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error or client disconnected:", err)
			break
		}
		
		log.Printf("Received: %s", message)

		// Echo the message back to the client
		response := fmt.Sprintf("Server Echo: %s", message)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve the WebSocket on the root path
	http.HandleFunc("/", handleWebSocket)

	fmt.Printf("Starting WebSocket server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
