package state

type State string

const (
	Above90   State = "ABOVE90"
	Under10   State = "UNDER10"
	CrossUp   State = "CROSS_UP"
	CrossDown State = "CROSS_DOWN"
	Neutral   State = "NEUTRAL"
)
