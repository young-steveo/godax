package market

// MinimumSize is a map of tickers to sizes
var MinimumSize = map[Ticker]Size{
	BTC: Size(`0.001`),
	BCH: Size(`0.01`),
	ETH: Size(`0.01`),
	LTC: Size(`0.1`),
	USD: Size(`10`),
	EUR: Size(`10`),
	GBP: Size(`10`),
}
