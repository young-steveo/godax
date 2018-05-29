package market

import "github.com/google/uuid"

// CompleteOrder is a command to remove orders to the book
type CompleteOrder struct {
	ID   uuid.UUID
	Resp chan *Order
}

// Type returns the command type
func (c *CompleteOrder) Type() CommandType {
	return Delete
}

// ClientID returns the client id
func (c *CompleteOrder) ClientID() uuid.UUID {
	return nilUUID
}

// ServerID returns the server id
func (c *CompleteOrder) ServerID() uuid.UUID {
	return c.ID
}

// Orders returns the order channel
func (c *CompleteOrder) Orders() chan *Order {
	return c.Resp
}
