package message

// Unsubscribe is a GDAX websocket unsubscribe message
var Unsubscribe = []byte(`{
	"type": "unsubscribe",
	"product_ids": ["LTC-USD"],
	"channels": ["heartbeat","ticker"]
}`)
