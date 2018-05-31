package market

import (
	"log"
	"strconv"
)

var epsilon = 0.00000001

// Price represents a float64 as a string (gdax API prefers strings).
type Price string

// Add two prices together
func (p1 Price) Add(p2 Price) (Price, error) {
	left, err := strconv.ParseFloat(string(p1), 64)
	if err != nil {
		return Price(`0.0`), err
	}
	right, err := strconv.ParseFloat(string(p2), 64)
	if err != nil {
		return Price(`0.0`), err
	}
	result := strconv.FormatFloat(left+right, 'f', 8, 64)
	return Price(result), nil
}

// Subtract p2 from this price
func (p1 Price) Subtract(p2 Price) (Price, error) {
	left, err := strconv.ParseFloat(string(p1), 64)
	if err != nil {
		log.Println(`Error parsing float:`, err.Error())
		return Price(`0.0`), err
	}
	right, err := strconv.ParseFloat(string(p2), 64)
	if err != nil {
		log.Println(`Error parsing float:`, err.Error())
		return Price(`0.0`), err
	}
	result := strconv.FormatFloat(left-right, 'f', 8, 64)
	return Price(result), nil
}

// Equals checks if two prices are the same with a margin of error for float precision.
func (p1 Price) Equals(p2 Price) (bool, error) {
	left, err := strconv.ParseFloat(string(p1), 64)
	if err != nil {
		log.Println(`Error parsing float:`, err.Error())
		return false, nil
	}
	right, err := strconv.ParseFloat(string(p2), 64)
	if err != nil {
		log.Println(`Error parsing float:`, err.Error())
		return false, nil
	}
	if (left-right) < epsilon && (right-left) < epsilon {
		return true, nil
	}
	return false, nil
}
