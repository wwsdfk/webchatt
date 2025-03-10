package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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
		return true // Разрешаем подключение откуда угодно
	},
}

var mutex = &sync.Mutex{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer conn.Close()

	// Получаем имя пользователя при подключении
	var initMsg Message
	err = conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Ошибка получения имени:", err)
		return
	}

	client := &Client{conn: conn, name: initMsg.Name}
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	// Чтение сообщений
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Ошибка чтения JSON:", err)
			mutex.Lock()
			delete(clients, client)
			mutex.Unlock()
			break
		}
		msg.Name = client.name // Добавляем имя в сообщение
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
				log.Println("Ошибка отправки JSON:", err)
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
	setupRoutes()

	go handleMessages()
	log.Fatal(http.ListenAndServe(":8083", nil))
}
