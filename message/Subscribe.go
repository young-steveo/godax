package message

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/young-steveo/godax/market"
)

type subscribeMessage struct {
	Typ        string             `json:"type"`
	ProductIDs []market.ProductID `json:"product_ids"`
	Channels   []string           `json:"channels"`
	Key        string             `json:"key,omitempty"`
	Passphrase string             `json:"passphrase,omitempty"`
	Timestamp  string             `json:"timestamp,omitempty"`
	Signature  string             `json:"signature,omitempty"`
}

// Subscribe is a GDAX websocket subscribe message
func Subscribe() []byte {
	key := os.Getenv(`GDAX_KEY`)
	passphrase := os.Getenv(`GDAX_PASSPHRASE`)
	stamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature, err := Signature([]byte(stamp+`GET/users/self/verify`), os.Getenv(`GDAX_SECRET`))
	if err != nil {
		log.Fatal("could not decode gdax secret", err)
	}

	message, _ := json.Marshal(subscribeMessage{
		`subscribe`,
		[]market.ProductID{market.ProductID{market.LTC, market.BTC}},
		[]string{`heartbeat`, `user`, `ticker`},
		key,
		passphrase,
		stamp,
		signature,
	})
	return message
}
