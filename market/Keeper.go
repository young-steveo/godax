package market

import (
	"log"

	"github.com/google/uuid"
)

// Keeper is a struct that manages state for a book via channels
type Keeper struct {
	book    MyBook
	adds    chan *Order
	reads   chan *GetByClientID
	updates chan *UpdateServerID
}

// GetByClientID is a command to get an order pointer by ClientID
type GetByClientID struct {
	ClientID uuid.UUID
	Orders   chan *Order
}

// UpdateServerID is a command to update the ServerID of an order that matches the ClientID
type UpdateServerID struct {
	ClientID uuid.UUID
	ServerID uuid.UUID
	Orders   chan *Order
}

// MakeKeeper is a constructor
func MakeKeeper(b MyBook, a chan *Order, r chan *GetByClientID, w chan *UpdateServerID) *Keeper {
	return &Keeper{b, a, r, w}
}

// Listen starts an infinite loop and consumes operations until the done channel tells it to quit.
func (k *Keeper) Listen(done chan bool) {
	defer close(done)
	for {
		select {
		case req := <-k.reads:
			req.Orders <- k.book.GetByClientID(req.ClientID)
		case req := <-k.updates:
			o := k.book.GetByClientID(req.ClientID)
			if o != nil {
				log.Printf(`Updating server ID of order %s -> %s`, o.ClientID, req.ServerID)
				o.ServerID = req.ServerID
			}
			req.Orders <- o
		case order := <-k.adds:
			log.Printf(`Adding order to local book %s`, order.ClientID)
			k.book = append(k.book, order)
		case <-done:
			return
		}
	}
}
