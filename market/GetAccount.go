package market

// GetAccount is a command to get an acocunt by Currency
type GetAccount struct {
	Ticker Ticker
	Resp   chan *Account
}

// Type returns the command type
func (c *GetAccount) Type() CommandType {
	return Read
}

// Currency returns the currency ticker
func (c *GetAccount) Currency() Ticker {
	return c.Ticker
}

// Accounts returns the account channel
func (c *GetAccount) Accounts() chan *Account {
	return c.Resp
}
