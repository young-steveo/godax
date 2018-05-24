package market

// CommandType is a numeric representation of command types
type CommandType int

// Command types
const (
	Create = iota
	Read
	Update
	Delete
)
