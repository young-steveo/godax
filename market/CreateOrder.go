package market

import "github.com/google/uuid"

// CreateOrder is a command to send orders to the book
type CreateOrder struct {
	Order *Order
	Resp  chan *Order
}

// Type returns the command type
func (c *CreateOrder) Type() CommandType {
	return Create
}

// ClientID returns the client id
func (c *CreateOrder) ClientID() uuid.UUID {
	return c.Order.ClientID
}

// ServerID returns the server id
func (c *CreateOrder) ServerID() uuid.UUID {
	return c.Order.ServerID
}

// Orders returns the order channel
func (c *CreateOrder) Orders() chan *Order {
	return c.Resp
}
