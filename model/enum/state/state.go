package state

type State string

const (
	Above90   State = "ABOVE90"
	Under10   State = "UNDER10"
	CrossUp   State = "CROSS_UP"
	CrossDown State = "CROSS_DOWN"
	Neutral   State = "NEUTRAL"
	Under0    State = "UNDER0"
	Above100  State = "ABOVE_100"
	Above1    State = "ABOVE_1"
)
