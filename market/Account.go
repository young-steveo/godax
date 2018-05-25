package market

import "github.com/google/uuid"

// Account is for a single currency
type Account struct {
	ID        uuid.UUID
	ProfileID uuid.UUID
	Currency  Ticker
	Balance   string
	Available string
	Hold      string
}
