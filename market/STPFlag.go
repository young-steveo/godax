package market

// STPFlag represents Self Trade Prevention Flags
type STPFlag string

// Self Trade Prevention Flags
const (
	DecreaseCancel = STPFlag(`dc`)
	CancelOldest   = STPFlag(`co`)
	CancelNewest   = STPFlag(`cn`)
	CancelBoth     = STPFlag(`cb`)
)
