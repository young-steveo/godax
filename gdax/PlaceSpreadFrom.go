package gdax

import (
	"log"
	"strconv"

	"github.com/google/uuid"
)

// PlaceSpreadFrom places two orders around the price of the indicated order.
// It will retry orders, moving them down or up until they are all successfull.
func PlaceSpreadFrom(orderID uuid.UUID) error {
	order := book.GetByServerID(orderID)
	p := order.Price
	price, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return err
	}
	PlaceSpread(order.ProductID, price)

	return nil
}
