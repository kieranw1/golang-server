package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// global variables
// map with key that is a pointer to a websocket - boolean value given for map
var clients = make(map[*websocket.Conn]bool)

// var for message channel
var broadcast = make(chan Message)

// create upgrader obj - takes http connections and upgrades to websockets
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// struct obj to hold messages
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// function to handle incoming websocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// use upgrader to upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// close connection when function returns
	defer ws.Close()

	// register new clients
	clients[ws] = true

	// infinite loop to pickup incoming messages
	for {
		var msg Message

		// read new messages as json and map to a message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			// if client has error, assume disconnection
			delete(clients, ws)
			break
		}
		// send new messages to broadcast channel
		broadcast <- msg
	}
}

// function to read broadcast channel and relay messages to clients
func handleMessages() {
	for {
		// get next message from broadcast channel
		msg := <-broadcast
		// send to every client connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				// assume client disconnection
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// application entrypoint
func main() {
	// create fileserver, tie to "/"
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// define websocket function route "/ws"
	http.HandleFunc("/ws", handleConnections)

	// listen for incoming chat messages - uses concurrency magic
	go handleMessages()

	// start server on localhost:8000, log any errors
	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
