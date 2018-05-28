package market

import (
	"fmt"

	"github.com/google/uuid"
)

// Account is for a single currency
type Account struct {
	ID        uuid.UUID
	ProfileID uuid.UUID
	Currency  Ticker
	Balance   string
	Available string
	Hold      string
}

func (a *Account) String() string {
	return fmt.Sprintf(`
		Account
		=======
		ID        %s
		ProfileID %s
		Currency  %s
		Balance   %s
		Available %s
		Hold      %s
	`, a.ID.String(), a.ProfileID.String(), a.Currency, a.Balance, a.Available, a.Hold)
}
