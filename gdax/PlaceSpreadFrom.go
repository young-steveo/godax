package gdax

import (
	"github.com/google/uuid"
)

// PlaceSpreadFrom places two orders around the price of the indicated order.
// It will retry orders, moving them down or up until they are all successfull.
func PlaceSpreadFrom(orderID uuid.UUID) error {
	order := book.GetByServerID(orderID)
	PlaceSpread(order.ProductID, order.Price)
	return nil
}
