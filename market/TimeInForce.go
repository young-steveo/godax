package market

// TimeInForce flags
type TimeInForce string

// Time in force flags
const (
	GoodTilCancel     = TimeInForce(`GTC`)
	GoodTilTime       = TimeInForce(`GTT`)
	ImmediateOrCancel = TimeInForce(`IOC`)
	ForkOrKill        = TimeInForce(`FOK`)
)
