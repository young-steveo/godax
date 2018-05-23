package market

import "github.com/google/uuid"

// MyBook represents all of my orders
type MyBook []*Order

// GetByClientID will find the order by ClientID
func (mb MyBook) GetByClientID(id uuid.UUID) *Order {
	for _, o := range mb {
		if o.ClientID == id {
			return o
		}
	}
	return nil
}
