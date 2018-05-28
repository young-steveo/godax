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
func Rebalance(pair market.ProductID, aCmd chan market.AccountCommand, oCmd chan market.OrderCommand) chan int {
	log.Printf(`Rebalancing %s pair`, pair)
	done := make(chan int)
	go func() {
		log.Printf(`Requesting %s product information`, pair)
		resp, err := Request(`GET`, fmt.Sprintf(`/products/%s/ticker`, pair), nil)
		if err != nil {
			log.Printf(`Error making request: %s`, err.Error())
			done <- -1
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf(`Error parsing body: %s`, err.Error())
			done <- -1
			return
		}

		fmt.Println(string(body))
		p, err := jsonparser.GetString(body, `price`)
		if err != nil {
			log.Printf(`Error parsing price: %s`, err.Error())
			done <- -1
			return
		}

		accounts := make(chan *market.Account)
		aCmd <- &market.GetAccount{Ticker: pair[0], Resp: accounts}
		leftAccount := <-accounts
		aCmd <- &market.GetAccount{Ticker: pair[1], Resp: accounts}
		rightAccount := <-accounts

		price, err := strconv.ParseFloat(p, 64)
		if err != nil {
			log.Printf(`Error converting price to float: %s`, err.Error())
			done <- -1
			return
		}

		leftBalance, err := strconv.ParseFloat(leftAccount.Available, 64)
		if err != nil {
			log.Printf(`Error converting price to float: %s`, err.Error())
			done <- -1
			return
		}

		rightBalance, err := strconv.ParseFloat(rightAccount.Available, 64)
		if err != nil {
			log.Printf(`Error converting price to float: %s`, err.Error())
			done <- -1
			return
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

			orders := make(chan *market.Order)
			defer close(orders)
			create := &market.CreateOrder{Order: o, Resp: orders}
			oCmd <- create
			create.Orders() <- o

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
				market.Price(strconv.FormatFloat(price-0.00001, 'f', 8, 64)),
				pair,
			)

			orders := make(chan *market.Order)
			defer close(orders)
			create := &market.CreateOrder{Order: o, Resp: orders}
			oCmd <- create
			create.Orders() <- o

			log.Printf(`Placing an order to %s %s %s for %s %s`, o.Side, o.Size, pair[0], o.Price, pair[1])
			Request(`POST`, `/orders`, o)
			orderMade = 1
		}

		done <- orderMade
	}()
	return done
}
