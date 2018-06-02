package market

import (
	"log"
	"strconv"
)

var epsilon = 0.00000001
var emptyPrice = Price(`0.0`)

// Price represents a float64 as a string (gdax API prefers strings).
type Price string

// FromFloat is a constructor
func FromFloat(val float64) Price {
	return Price(strconv.FormatFloat(val, 'f', 8, 64))
}

// Add two prices together
func (p1 Price) Add(p2 Price) (Price, error) {
	left, err := p1.float()
	if err != nil {
		return emptyPrice, err
	}
	right, err := p2.float()
	if err != nil {
		return emptyPrice, err
	}
	return FromFloat(left + right), nil
}

// Subtract p2 from this price
func (p1 Price) Subtract(p2 Price) (Price, error) {
	left, err := p1.float()
	if err != nil {
		return emptyPrice, err
	}
	right, err := p2.float()
	if err != nil {
		return emptyPrice, err
	}
	return FromFloat(left - right), nil
}

// Equals checks if two prices are the same with a margin of error for float precision.
func (p1 Price) Equals(p2 Price) (bool, error) {
	left, err := p1.float()
	if err != nil {
		return false, err
	}
	right, err := p2.float()
	if err != nil {
		return false, err
	}
	if (left-right) < epsilon && (right-left) < epsilon {
		return true, nil
	}
	return false, nil
}

// Less will check if p1 is less than p2
func (p1 Price) Less(p2 Price) (bool, error) {
	left, err := p1.float()
	if err != nil {
		return false, err
	}
	right, err := p2.float()
	if err != nil {
		return false, err
	}
	return left < right, nil
}

func (p1 Price) float() (float64, error) {
	p, err := strconv.ParseFloat(string(p1), 64)
	if err != nil {
		log.Println(`Error parsing float:`, err.Error())
		return 0.0, err
	}
	return p, nil
}
