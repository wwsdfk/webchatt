package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"chat-app/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Client struct {
	conn *websocket.Conn
	name string
}

var clients = make(map[*Client]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mutex = &sync.Mutex{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении соединения: %v\n", err)
		return
	}
	defer conn.Close()

	var initMsg Message
	err = conn.ReadJSON(&initMsg)
	if err != nil {
		log.Printf("Ошибка получения имени: %v\n", err)
		return
	}

	client := &Client{conn: conn, name: initMsg.Name}
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Ошибка чтения JSON: %v\n", err)
			mutex.Lock()
			delete(clients, client)
			mutex.Unlock()
			break
		}
		msg.Name = client.name

		_, err = database.DB.Exec(r.Context(), "INSERT INTO messages (name, content) VALUES ($1, $2)", msg.Name, msg.Content)
		if err != nil {
			log.Printf("Ошибка сохранения в базу данных: %v\n", err)
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка отправки JSON: %v\n", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleConnections)
	http.Handle("/", r)
}

func main() {
	fmt.Println("Сервер запущен на порту :8083")
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.DB.Close()

	setupRoutes()
	go handleMessages()
	log.Fatal(http.ListenAndServe(":8083", nil))
}
