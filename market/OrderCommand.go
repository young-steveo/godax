package market

import "github.com/google/uuid"

// OrderCommand defines a command for the Keeper
type OrderCommand interface {
	Type() CommandType
	ClientID() uuid.UUID
	ServerID() uuid.UUID
	Orders() chan *Order
}
