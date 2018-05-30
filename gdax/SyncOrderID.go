package gdax

import (
	"log"

	"github.com/google/uuid"
)

// SyncOrderID updates an order's ServerID by ClientID
func SyncOrderID(clientID uuid.UUID, orderID uuid.UUID) {
	log.Printf(`Synchronizing order server id with local client id. ClientID: %s ServerID: %s`, clientID, orderID)
	book.SyncOrderID(clientID, orderID)
}
