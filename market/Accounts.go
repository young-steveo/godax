package market

// Accounts represents all of my orders
type Accounts []*Account

// GetByCurrency will find the account by Currency
func (accts Accounts) GetByCurrency(currency Ticker) *Account {
	for _, a := range accts {
		if a.Currency == currency {
			return a
		}
	}
	return nil
}
