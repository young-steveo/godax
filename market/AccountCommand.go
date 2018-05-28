package market

// AccountCommand defines a command for the Keeper
type AccountCommand interface {
	Type() CommandType
	Currency() Ticker
	Accounts() chan *Account
}
