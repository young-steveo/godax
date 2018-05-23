package gdax

import (
	"github.com/young-steveo/godax/market"
)

// PlaceOrder will send an order struct over to GDAX
func PlaceOrder(o *market.Order) {
	request(`POST`, `/orders`, o)
}
