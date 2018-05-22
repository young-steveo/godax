package gdax

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Connect to the websocket
func Connect() error {
	u := url.URL{Scheme: "wss", Host: "ws-feed-public.sandbox.gdax.com"}
	//u := url.URL{Scheme: "wss", Host: "ws-feed.gdax.com"}

	log.Printf("Connecting to %s", u.String())
	var err error
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		log.Printf("Connected")
	}
	return err
}
