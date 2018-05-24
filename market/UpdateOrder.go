package market

import "github.com/google/uuid"

// UpdateOrder is a command to update the ServerID of an order that matches the ClientID
type UpdateOrder struct {
	ClientIDValue uuid.UUID
	ServerIDValue uuid.UUID
	Resp          chan *Order
}

// Type returns the command type
func (c *UpdateOrder) Type() CommandType {
	return Update
}

// ClientID returns the client id
func (c *UpdateOrder) ClientID() uuid.UUID {
	return c.ClientIDValue
}

// ServerID returns the server id
func (c *UpdateOrder) ServerID() uuid.UUID {
	return c.ServerIDValue
}

// Orders returns the order channel
func (c *UpdateOrder) Orders() chan *Order {
	return c.Resp
}
