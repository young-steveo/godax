package gdax

// ReadMessage a GDAX message from the websocket
func ReadMessage() (int, []byte, error) {
	return ws.ReadMessage()
}
