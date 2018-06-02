package gdax

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/young-steveo/godax/market"
)

// Rebalance will try to use limit orders to make sure both assets of a pair are balanced.
func Rebalance(pair market.ProductID) (int, float64, error) {
	log.Printf(`Rebalancing %s pair`, pair)
	log.Printf(`Requesting %s product information`, pair)
	resp, err := Request(`GET`, fmt.Sprintf(`/products/%s/ticker`, pair), nil)
	if err != nil {
		log.Printf(`Error making request: %s`, err.Error())
		return -1, 0.0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(`Error parsing body: %s`, err.Error())
		return -1, 0.0, err
	}
	p, err := jsonparser.GetString(body, `price`)
	if err != nil {
		log.Printf(`Error parsing price: %s`, err.Error())
		return -1, 0.0, err
	}
	leftAccount := accounts.GetByCurrency(pair[0])
	rightAccount := accounts.GetByCurrency(pair[1])

	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return -1, 0.0, err
	}

	leftBalance, err := strconv.ParseFloat(leftAccount.Balance, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return -1, 0.0, err
	}

	rightBalance, err := strconv.ParseFloat(rightAccount.Balance, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return -1, 0.0, err
	}

	allRight := rightBalance + (leftBalance * price)
	newRightBal := allRight / 2
	newLeftBal := newRightBal / price

	log.Printf(`Current %s Balance: %f`, pair[0], leftBalance)
	log.Printf(`Current %s Balance: %f`, pair[1], rightBalance)

	log.Printf(`Target %s Balance: %f`, pair[0], newLeftBal)

	orderMade := 0

	if newLeftBal > leftBalance*1.2 {
		log.Printf(
			`Target %s balance of %f is greater than current balance %f * 0.8 (%f)`,
			pair[0],
			newLeftBal,
			leftBalance,
			leftBalance*1.2,
		)

		amtToBuy := newLeftBal - leftBalance
		o := market.MakeOrder(
			market.Side("buy"),
			market.Size(strconv.FormatFloat(amtToBuy, 'f', 8, 64)),

			// fixme (should not be hardcoded, should get increments from gdax)
			market.Price(strconv.FormatFloat(price-0.00001, 'f', 8, 64)),
			pair,
		)

		book.AddOrder(o)

		log.Printf(`Placing an order to %s %s %s for %s %s`, o.Side, o.Size, pair[0], o.Price, pair[1])
		Request(`POST`, `/orders`, o)

		orderMade = 1

	} else if newLeftBal < leftBalance*0.8 {
		log.Printf(
			`Target %s balance of %f is less than current balance %f * 0.8 (%f)`,
			pair[0],
			newLeftBal,
			leftBalance,
			leftBalance*0.8,
		)

		amtToSell := leftBalance - newLeftBal
		o := market.MakeOrder(
			market.Side("sell"),
			market.Size(strconv.FormatFloat(amtToSell, 'f', 8, 64)),

			// fixme (should not be hardcoded, should get increments from gdax)
			market.Price(strconv.FormatFloat(price+0.00001, 'f', 8, 64)),
			pair,
		)

		book.AddOrder(o)

		log.Printf(`Placing an order to %s %s %s for %s %s`, o.Side, o.Size, pair[0], o.Price, pair[1])
		Request(`POST`, `/orders`, o)
		orderMade = 1
	}

	return orderMade, allRight, nil
}
