package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	log.Println("Client connected to WebSocket!")

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

// Default HTTP Endpoint
func handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP request on %s", r.URL.Path)
	w.Write([]byte("Go App is running perfectly!"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ---------------------------------------------------------
	// Background loop to constantly generate logs
	// This makes it easy to test your Real-Time Log Streaming!
	// ---------------------------------------------------------
	go func() {
		for {
			log.Println("[INFO] Application heartbeat: Running smoothly...")

			log.Println("[INFO] github auto webhook...")
			time.Sleep(3 * time.Second)
		}
	}()

	// Mount the routes
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", handleDefault)

	log.Printf("Starting HTTP & WebSocket server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
