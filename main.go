package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MaxRubel/WebsocketsGo/db"
	"github.com/MaxRubel/WebsocketsGo/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
	
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("unable to print incoming message")
			return
		}

		parsedMessage, err := models.ParseJSONMessage(message)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		db.AddMessageToDb(parsedMessage)

		fmt.Println("received message:", parsedMessage)

		messages, err := db.GetAllMessages()

		if err != nil {
			log.Println(err)
		}

		messageF, err := json.Marshal(messages)

		if err != nil {
			log.Println(err)
		}

        err = conn.WriteMessage(messageType, messageF)
        if err != nil {
            log.Println("Write error:", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", wsHandler)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        log.Fatal("ListenAndServe error:", err)
    }

}
