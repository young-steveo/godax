package market

import (
	"sync"
)

// Accounts represents all of my accounts
type Accounts struct {
	List  []*Account
	Mutex *sync.Mutex
}

// GetAccounts is a constructor
func GetAccounts() *Accounts {
	return &Accounts{make([]*Account, 0), &sync.Mutex{}}
}

// GetByCurrency will find the account by Currency
func (accts *Accounts) GetByCurrency(currency Ticker) *Account {
	accts.Mutex.Lock()
	for _, a := range accts.List {
		if a.Currency == currency {
			accts.Mutex.Unlock()
			return a
		}
	}
	accts.Mutex.Unlock()
	return nil
}

// AddAccount adds an account to the list
func (accts *Accounts) AddAccount(acct *Account) {
	accts.Mutex.Lock()
	accts.List = append(accts.List, acct)
	accts.Mutex.Unlock()
}
