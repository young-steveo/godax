package message

// Unsubscribe is a GDAX websocket unsubscribe message
func Unsubscribe() []byte {
	return []byte(`{
		"type": "unsubscribe",
		"product_ids": ["LTC-USD"],
		"channels": ["heartbeat","user","ticker"]
	}`)
}
