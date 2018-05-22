package gdax

// Close the websocket connection to GDAX
func Close() {
	ws.Close()
}
