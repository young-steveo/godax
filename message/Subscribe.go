package message

// Subscribe is a GDAX websocket subscribe message
var Subscribe = []byte(`{
	"type": "subscribe",
	"product_ids": ["LTC-USD"],
	"channels": [
		"heartbeat",
		{
			"name":"ticker",
			"product_ids": ["LTC-USD"]
		}
	]
}`)
