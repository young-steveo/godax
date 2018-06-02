package market

// ProductIncrement is a map of an Array of Tickers to a minimum Price
var ProductIncrement = map[ProductID]Price{
	ProductID{BTC, USD}: Price(`0.01`),
	ProductID{BTC, EUR}: Price(`0.01`),
	ProductID{BTC, GBP}: Price(`0.01`),
	ProductID{ETH, USD}: Price(`0.01`),
	ProductID{ETH, EUR}: Price(`0.01`),
	ProductID{ETH, BTC}: Price(`0.00001`),
	ProductID{LTC, USD}: Price(`0.01`),
	ProductID{LTC, EUR}: Price(`0.01`),
	ProductID{LTC, BTC}: Price(`0.00001`),
	ProductID{BCH, USD}: Price(`0.01`),
	ProductID{BCH, EUR}: Price(`0.01`),
	ProductID{BCH, BTC}: Price(`0.00001`),
}
