package gdax

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/young-steveo/godax/market"
)

// GetBalance will return a float64 of the right side of the ProductID if all the left were to be converted.
func GetBalance(pair market.ProductID) (float64, market.Price, error) {
	log.Printf(`Requesting %s product information`, pair)
	resp, err := Request(`GET`, fmt.Sprintf(`/products/%s/ticker`, pair), nil)
	if err != nil {
		log.Printf(`Error making request: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(`Error parsing body: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}
	p, err := jsonparser.GetString(body, `price`)
	if err != nil {
		log.Printf(`Error parsing price: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}
	leftAccount := accounts.GetByCurrency(pair[0])
	rightAccount := accounts.GetByCurrency(pair[1])

	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}

	leftBalance, err := strconv.ParseFloat(leftAccount.Balance, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}

	rightBalance, err := strconv.ParseFloat(rightAccount.Balance, 64)
	if err != nil {
		log.Printf(`Error converting price to float: %s`, err.Error())
		return 0.0, market.Price(`0.0`), err
	}

	return rightBalance + (leftBalance * price), market.Price(p), nil
}
