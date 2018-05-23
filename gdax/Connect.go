package gdax

import (
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// Connect to the websocket
func Connect() error {
	var host string
	switch os.Getenv(`GDAX_MODE`) {
	case `prod`:
		host = os.Getenv(`GDAX_FEED_HOST`)
	case `dev`:
		host = os.Getenv(`GDAX_SANDBOX_FEED_HOST`)
	default:
		return errors.New(`env variable GDAX_MODE not set or invalid`)
	}
	u := url.URL{Scheme: "wss", Host: host}

	log.Printf("Connecting to %s", u.String())
	var err error
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		log.Printf("Connected")
	}
	return err
}
