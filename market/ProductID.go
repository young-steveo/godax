package market

import (
	"fmt"
)

// ProductID is a pair of tickers
type ProductID [2]Ticker

// GetProductID will return a ProductID pair of tickers.
func GetProductID(left string, right string) ProductID {
	return ProductID{Ticker(left), Ticker(right)}
}

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (pid ProductID) MarshalText() (text []byte, err error) {
	return []byte(pid.String()), nil
}

func (pid ProductID) String() string {
	return fmt.Sprintf(`%s-%s`, pid[0], pid[1])
}
