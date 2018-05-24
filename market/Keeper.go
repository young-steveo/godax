package market

import (
	"log"
)

// Keeper is a struct that manages state for a book via channels
type Keeper struct {
	book     MyBook
	commands chan OrderCommand
}

// MakeKeeper is a constructor
func MakeKeeper(book MyBook, commands chan OrderCommand) *Keeper {
	return &Keeper{book, commands}
}

// Listen starts an infinite loop and consumes operations until the done channel tells it to quit.
func (k *Keeper) Listen(done chan bool) {
	defer close(done)
	for {
		select {
		case command := <-k.commands:
			k.handleCommand(command)
		case <-done:
			return
		}
	}
}

func (k *Keeper) handleCommand(command OrderCommand) {
	switch command.Type() {
	case Read:
		command.Orders() <- k.book.GetByClientID(command.ClientID())
	case Update:
		o := k.book.GetByClientID(command.ClientID())
		if o != nil {
			serverID := command.ServerID()
			log.Printf(`Updating server ID of order %s -> %s`, o.ClientID, serverID)
			o.ServerID = serverID
		}
		command.Orders() <- o
	case Create:
		order := <-command.Orders()
		log.Printf(`Adding order to local book %s`, order.ClientID)
		k.book = append(k.book, order)
	}
}
