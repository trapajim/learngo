package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var messages = make(chan Message)
var upgrader = websocket.Upgrader{}

func main() {
	clientFs := http.FileServer(http.Dir("./client/build"))
	http.Handle("/", clientFs)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(responseWriter http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(responseWriter, request, nil)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		messages <- msg
	}
}

func handleMessages() {
	for {
		msg := <-messages

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
