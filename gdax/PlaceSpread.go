package gdax

import (
	"io/ioutil"
	"log"

	"github.com/buger/jsonparser"
	"github.com/young-steveo/godax/market"
)

// PlaceSpread will make two orders
func PlaceSpread(pid market.ProductID, price market.Price) {
	log.Printf(`Placing Spread order for %s around %f`, pid, price)
	increment(pid, price, 0)
	decrement(pid, price, 0)
}

func increment(pid market.ProductID, price market.Price, times int) bool {
	price, err := price.Add(market.Price(`0.00001`))
	if err != nil {
		log.Println(`Could not add prices for increment: ` + err.Error())
	}
	old := book.GetByPrice(price)
	if old != nil {
		log.Printf(`Existing order for price %s`, price)
		return true
	}
	sell := market.MakeOrder(market.Side(`sell`), market.Size(`0.1`), price, pid)

	book.AddOrder(sell)

	log.Printf(`%s %s %s for %s %s`, sell.Side, sell.Size, pid[0], sell.Price, pid[1])
	res, err := Request(`POST`, `/orders`, sell)
	if err != nil {
		log.Println(`Could not make sell request: ` + err.Error())
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(`Could not parse body: ` + err.Error())
		return false
	}
	status, err := jsonparser.GetUnsafeString(body, `status`)
	if err != nil {
		log.Println(`Could not get status: ` + err.Error())
		return false
	}
	if status == `rejected` && times < 10 {
		log.Printf(`Sell order rejected, increasing price to %s`, price)
		return increment(pid, price, times+1)
	}
	return times < 10
}

func decrement(pid market.ProductID, price market.Price, times int) bool {
	price, err := price.Subtract(market.Price(`0.00001`))
	if err != nil {
		log.Println(`Could not subtract prices for increment: ` + err.Error())
	}
	old := book.GetByPrice(price)
	if old != nil {
		log.Printf(`Existing order for price %s`, price)
		return true
	}
	buy := market.MakeOrder(market.Side(`buy`), market.Size(`0.1`), price, pid)

	book.AddOrder(buy)

	log.Printf(`%s %s %s for %s %s`, buy.Side, buy.Size, pid[0], buy.Price, pid[1])
	res, err := Request(`POST`, `/orders`, buy)
	if err != nil {
		log.Println(`Could not make buy request: ` + err.Error())
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(`Could not parse body: ` + err.Error())
		return false
	}
	status, err := jsonparser.GetUnsafeString(body, `status`)
	if err != nil {
		log.Println(`Could not get status: ` + err.Error())
		return false
	}
	if status == `rejected` && times < 10 {
		log.Printf(`Buy order rejected, decreasing price to %s`, price)
		// increment(pid, price, 9)
		return decrement(pid, price, times+1)
	}
	return times < 10
}
