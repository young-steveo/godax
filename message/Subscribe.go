package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

type subscribeMessage struct {
	Typ        string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
	Key        string   `json:"key,omitempty"`
	Passphrase string   `json:"passphrase,omitempty"`
	Timestamp  string   `json:"timestamp,omitempty"`
	Signature  string   `json:"signature,omitempty"`
}

// Subscribe is a GDAX websocket subscribe message
func Subscribe() []byte {
	key := os.Getenv(`GDAX_KEY`)
	passphrase := os.Getenv(`GDAX_PASSPHRASE`)
	secret, err := base64.StdEncoding.DecodeString(os.Getenv(`GDAX_SECRET`))
	if err != nil {
		log.Fatal("could not decode gdax secret", err)
	}

	stamp := strconv.FormatInt(time.Now().Unix(), 10)

	signature := hmac.New(sha256.New, secret)

	signature.Write([]byte(stamp + `GET/users/self/verify`))

	message, _ := json.Marshal(subscribeMessage{
		`subscribe`,
		[]string{`LTC-USD`},
		[]string{`heartbeat`, `user`, `ticker`},
		key,
		passphrase,
		stamp,
		base64.StdEncoding.EncodeToString(signature.Sum(nil)),
	})
	return message
}
