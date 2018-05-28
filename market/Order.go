package market

import (
	"github.com/google/uuid"
)

// Size of the order.
type Size string

// Price of the order.
type Price string

// Side of the order "buy" or "sell"
type Side string

// Order is a json serializeable order struct
type Order struct {
	ServerID    uuid.UUID   `json:"-"`
	ClientID    uuid.UUID   `json:"client_oid"`
	Type        string      `json:"type"`
	Side        Side        `json:"side"`
	ProductID   ProductID   `json:"product_id"`
	Stp         STPFlag     `json:"stp"`
	Size        Size        `json:"size"`
	Price       Price       `json:"price"`
	TimeInForce TimeInForce `json:"time_in_force"`
	PostOnly    bool        `json:"post_only"`
}

// MakeOrder makes a new order
func MakeOrder(side Side, size Size, price Price, pair ProductID) *Order {
	clientID := uuid.New()
	return &Order{
		ClientID:    clientID,
		Type:        `limit`,
		Side:        side,
		ProductID:   pair,
		Stp:         CancelOldest,
		Size:        size,
		Price:       price,
		TimeInForce: GoodTilCancel,
		PostOnly:    true,
	}
}
