package market

// AddAccount is a command to send add accounts to the book
type AddAccount struct {
	Ticker Ticker
	Resp   chan *Account
}

// Type returns the command type
func (c *AddAccount) Type() CommandType {
	return Create
}

// Currency returns the currency ticker
func (c *AddAccount) Currency() Ticker {
	return c.Ticker
}

// Accounts returns the account channel
func (c *AddAccount) Accounts() chan *Account {
	return c.Resp
}
