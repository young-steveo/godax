package market

import (
	"sync"

	"github.com/google/uuid"
)

// Book represents all of my orders
type Book struct {
	List  []*Order
	Mutex *sync.Mutex
}

// GetBook is a constructor
func GetBook() *Book {
	return &Book{make([]*Order, 0), &sync.Mutex{}}
}

// GetByClientID will find the order by ClientID
func (book *Book) GetByClientID(id uuid.UUID) *Order {
	book.Mutex.Lock()
	for _, o := range book.List {
		if o.ClientID == id {
			book.Mutex.Unlock()
			return o
		}
	}
	book.Mutex.Unlock()
	return nil
}

// AddOrder adds an order to the list
func (book *Book) AddOrder(order *Order) {
	book.Mutex.Lock()
	book.List = append(book.List, order)
	book.Mutex.Unlock()
}
