package market

import (
	"fmt"
)

// ProductID is a pair of tickers
type ProductID []Ticker

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (pid ProductID) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf(`%s-%s`, pid[0], pid[1])), nil
}
