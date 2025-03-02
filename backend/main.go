package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Client struct {
	conn *websocket.Conn
}

var clients = make(map[*Client]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	client := &Client{conn: conn}
	clients[client] = true
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Ошибка чтения JSON:", err)
			delete(clients, client)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Println("Ошибка отправки JSON:", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
	}
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleConnections)
	http.Handle("/", r)
}

func main() {
	fmt.Println("Сервер запущен на порту :8080")
	setupRoutes()
	go handleMessages()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
