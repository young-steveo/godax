package market

import (
	"github.com/google/uuid"
)

// Order is a json serializeable order struct
type Order struct {
	ServerID    uuid.UUID   `json:"-"`
	ClientID    uuid.UUID   `json:"client_oid"`
	Typ         string      `json:"type"`
	Side        string      `json:"side"`
	ProductID   string      `json:"product_id"`
	Stp         STPFlag     `json:"stp"`
	Size        string      `json:"size"`
	Price       string      `json:"price"`
	TimeInForce TimeInForce `json:"time_in_force"`
	PostOnly    bool        `json:"post_only"`
}

// MakeOrder makes a new order
func MakeOrder(side string, size string, price string) *Order {
	clientID := uuid.New()
	return &Order{
		ClientID:    clientID,
		Typ:         `limit`,
		Side:        side,
		ProductID:   `LTC-USD`,
		Stp:         CancelOldest,
		Size:        size,
		Price:       price,
		TimeInForce: GoodTilCancel,
		PostOnly:    true,
	}
}
