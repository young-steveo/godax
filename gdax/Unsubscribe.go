package gdax

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/young-steveo/godax/message"
)

// Unsubscribe from websocket
func Unsubscribe() {
	if ws != nil {
		log.Printf("Unsubscribing to channels")
		err := ws.WriteMessage(websocket.TextMessage, []byte(message.Unsubscribe))
		if err != nil {
			log.Printf("Could not unsubscribe: %s\n", err.Error())
		}
	}
}
