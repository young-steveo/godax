package gdax

import (
	"io/ioutil"

	"github.com/google/uuid"

	"github.com/buger/jsonparser"
	"github.com/young-steveo/godax/market"
)

// GetAccounts will return a slice of market.Account pointers
func GetAccounts() ([]*market.Account, error) {
	var accounts []*market.Account

	response, err := Request(`GET`, `/accounts`, nil)
	if err != nil {
		return accounts, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return accounts, err
	}
	var broken error
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if broken != nil {
			return
		}
		idBytes, _, _, err := jsonparser.Get(value, `id`)
		if err != nil {
			broken = err
			return
		}
		id, err := uuid.ParseBytes(idBytes)
		if err != nil {
			broken = err
			return
		}
		PIDBytes, _, _, err := jsonparser.Get(value, `id`)
		if err != nil {
			broken = err
			return
		}
		ProfileID, err := uuid.ParseBytes(PIDBytes)
		if err != nil {
			broken = err
			return
		}
		currency, err := jsonparser.GetString(value, `currency`)
		if err != nil {
			broken = err
			return
		}
		balance, err := jsonparser.GetString(value, `balance`)
		if err != nil {
			broken = err
			return
		}
		available, err := jsonparser.GetString(value, `available`)
		if err != nil {
			broken = err
			return
		}
		hold, err := jsonparser.GetString(value, `hold`)
		if err != nil {
			broken = err
			return
		}

		accounts = append(accounts, &market.Account{
			ID:        id,
			ProfileID: ProfileID,
			Currency:  market.Ticker(currency),
			Balance:   balance,
			Available: available,
			Hold:      hold,
		})
	})
	return accounts, broken
}
