package market

import (
	"log"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

var epsilon = 0.0000001

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

// GetByServerID will find the order by ClientID
func (book *Book) GetByServerID(id uuid.UUID) *Order {
	book.Mutex.Lock()
	for _, o := range book.List {
		if o.ServerID == id {
			book.Mutex.Unlock()
			return o
		}
	}
	book.Mutex.Unlock()
	return nil
}

// GetByPrice will find the order by Price and ProductID
func (book *Book) GetByPrice(price float64) *Order {
	book.Mutex.Lock()
	for _, o := range book.List {
		oPrice, err := strconv.ParseFloat(string(o.Price), 64)
		if err != nil {
			book.Mutex.Unlock()
			log.Printf(`Error converting price to float: %s`, err.Error())
			return nil
		}
		if (price-oPrice) < epsilon && (oPrice-price) < epsilon {
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

// SyncOrderID updates an order's ServerID by ClientID
func (book *Book) SyncOrderID(clientID uuid.UUID, orderID uuid.UUID) {
	order := book.GetByClientID(clientID)
	book.Mutex.Lock()
	order.ServerID = orderID
	book.Mutex.Unlock()
}

// RemoveOrder removes an order by serverID
func (book *Book) RemoveOrder(orderID uuid.UUID) {
	var index int
	var o *Order
	book.Mutex.Lock()
	for index, o = range book.List {
		if o.ServerID == orderID {
			book.List = append(book.List[:index], book.List[index+1:]...)
			book.Mutex.Unlock()
			return
		}
	}
	book.Mutex.Unlock()
}
