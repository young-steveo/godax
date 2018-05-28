package market

import (
	"log"
)

// Keeper is a struct that manages state for a book via channels
type Keeper struct {
	accounts        Accounts
	book            MyBook
	orderCommands   chan OrderCommand
	accountCommands chan AccountCommand
}

// MakeKeeper is a constructor
func MakeKeeper(oCommands chan OrderCommand, aCommands chan AccountCommand) *Keeper {
	book := make(MyBook, 0)
	accounts := make(Accounts, 0)
	return &Keeper{accounts, book, oCommands, aCommands}
}

// Listen starts an infinite loop and consumes operations until the done channel tells it to quit.
func (k *Keeper) Listen(done chan bool) {
	defer close(done)
	for {
		select {
		case command := <-k.orderCommands:
			k.handleOrderCommand(command)
		case command := <-k.accountCommands:
			k.handleAccountCommand(command)
		case <-done:
			return
		}
	}
}

func (k *Keeper) handleOrderCommand(command OrderCommand) {
	switch command.Type() {
	case Create:
		order := <-command.Orders()
		log.Printf(`Adding order to local book %s`, order.ClientID)
		k.book = append(k.book, order)
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
	}
}

func (k *Keeper) handleAccountCommand(command AccountCommand) {
	switch command.Type() {
	case Create:
		account := <-command.Accounts()
		log.Printf(`Adding account to local accounts %s`, account.Currency)
		k.accounts = append(k.accounts, account)
		command.Accounts() <- account // send it back so we know we're done.
	}
}
