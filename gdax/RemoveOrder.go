package gdax

import (
	"log"

	"github.com/google/uuid"
)

// RemoveOrder kills an order on the book
func RemoveOrder(orderID uuid.UUID) {
	log.Printf(`Removing order %s from the book.`, orderID)
	book.RemoveOrder(orderID)
}
