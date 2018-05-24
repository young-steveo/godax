package market

import "github.com/google/uuid"

// GetByClientID is a command to get an order pointer by ClientID
type GetByClientID struct {
	ID   uuid.UUID
	Resp chan *Order
}

// Type returns the command type
func (c *GetByClientID) Type() CommandType {
	return Read
}

// ClientID returns the client id
func (c *GetByClientID) ClientID() uuid.UUID {
	return c.ID
}

// ServerID returns the server id
func (c *GetByClientID) ServerID() uuid.UUID {
	return nilUUID
}

// Orders returns the order channel
func (c *GetByClientID) Orders() chan *Order {
	return c.Resp
}
