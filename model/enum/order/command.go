package order

type Command string

const (
	LONG  Command = "LONG"
	SHORT Command = "SHORT"
	WAIT  Command = "WAIT"
	CLOSE Command = "CLOSE"
)
