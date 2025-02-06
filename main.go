package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var mu = &sync.Mutex{}
var rooms = make(map[string]*Room)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var chat = &Chat{}

func main() {
	http.HandleFunc("/ws", serveWs)

	http.HandleFunc("/", http.FileServer(http.Dir("./public")).ServeHTTP)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, msgReader, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}

		if messageType != websocket.TextMessage {
			log.Println("invalid message type")
			continue
		}

		var msg GenericMessage
		err = json.NewDecoder(msgReader).Decode(&msg)
		if err != nil {
			log.Printf("failed to decode message: %s", err)
			continue
		}

		switch msg.Type {
		case CreateRoomCommand:
			var body CreateRoomMessageBody
			err = json.Unmarshal(msg.Body, &body)
			if err != nil {
				log.Printf("failed to unmarshal create room message body: %s", err)
				continue
			}

			room, err := chat.CreateRoom(ctx, body.RoomName)
			if err != nil {
				log.Printf("failed to create room %s: %s", body.RoomName, err)
				continue
			}

		}
	}
}
