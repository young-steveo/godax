package gdax

import (
	"log"

	"github.com/young-steveo/godax/market"

	"github.com/gorilla/websocket"
	"github.com/young-steveo/godax/message"
)

// Subscribe to GDAX messages from channels
func Subscribe(productID market.ProductID) {
	log.Println(`Subscribing to channels`)
	err := ws.WriteMessage(websocket.TextMessage, message.Subscribe(productID))
	if err != nil {
		log.Fatal("Could not subscribe: ", err)
		return
	}
	log.Println(`Subscription sent`)
}
