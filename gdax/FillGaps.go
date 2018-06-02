package gdax

import (
	"io/ioutil"
	"log"

	"github.com/google/uuid"

	"github.com/buger/jsonparser"
	"github.com/young-steveo/godax/market"
)

// FillGaps will look at the order book and place orders for anything missing between
// the lowest buy and the highest sell.
func FillGaps(pid market.ProductID) {
	prices := book.GapPrices(pid)
	size := market.MinimumSize[pid[0]]
	for price, side := range prices {
		order := market.MakeOrder(side, size, price, pid)

		book.AddOrder(order)

		log.Printf(`%s %s %s for %s %s`, side, size, pid[0], price, pid[1])
		res, err := Request(`POST`, `/orders`, order)
		if err != nil {
			book.RemoveOrder(order.ClientID)
			log.Println(`Could not make sell request: ` + err.Error())
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			book.RemoveOrder(order.ClientID)
			log.Println(`Could not parse body: ` + err.Error())
			return
		}
		status, err := jsonparser.GetUnsafeString(body, `status`)
		if err != nil {
			book.RemoveOrder(order.ClientID)
			log.Println(`Could not get status: ` + err.Error())
			return
		}
		if status == `rejected` {
			book.RemoveOrder(order.ClientID)
			reason, err := jsonparser.GetUnsafeString(body, `reject_reason`)
			if err != nil {
				log.Println(`Could not get reject_reason: ` + err.Error())
				return
			}
			if reason == `post only` {
				// might be wrong side?  Hacky but this should work...
				if string(order.Side) == `buy` {
					order.Side = market.Side(`sell`)
				} else {
					order.Side = market.Side(`buy`)
				}
				order.ClientID = uuid.New()
				book.AddOrder(order)
				log.Printf(`%s %s %s for %s %s`, side, size, pid[0], price, pid[1])
				res, err := Request(`POST`, `/orders`, order)
				if err != nil {
					book.RemoveOrder(order.ClientID)
					log.Println(`Could not make sell request: ` + err.Error())
					return
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					book.RemoveOrder(order.ClientID)
					log.Println(`Could not parse body: ` + err.Error())
					return
				}
				status, err := jsonparser.GetUnsafeString(body, `status`)
				if err != nil {
					book.RemoveOrder(order.ClientID)
					log.Println(`Could not get status: ` + err.Error())
					return
				}
				if status == `rejected` {
					book.RemoveOrder(order.ClientID)
				}
			}
			log.Printf(`Fill Gap order Rejected: %s`, string(body))
		}
	}
}
