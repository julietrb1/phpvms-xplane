package models

type PirepStatus string

const (
	PIREPStatusInitiated    PirepStatus = "INI"
	PIREPStatusBoarding     PirepStatus = "BST"
	PIREPStatusTaxiing      PirepStatus = "TXI"
	PIREPStatusTakeOff      PirepStatus = "TOF"
	PIREPStatusTakeOffClimb PirepStatus = "TKO"
	PIREPStatusEnRoute      PirepStatus = "ENR"
	PIREPStatusTopOfDescent PirepStatus = "TEN"
	PIREPStatusLanding      PirepStatus = "LDG"
	PIREPStatusLanded       PirepStatus = "LAN"
	PIREPStatusArrived      PirepStatus = "ARR"
	PIREPStatusDeboarding   PirepStatus = "DX"
	PIREPStatusPostShutdown PirepStatus = "PSD"
)

func ValidateStatus(status string) bool {
	switch PirepStatus(status) {
	case PIREPStatusInitiated, PIREPStatusBoarding,
		PIREPStatusTaxiing, PIREPStatusTakeOff, PIREPStatusTakeOffClimb,
		PIREPStatusEnRoute, PIREPStatusTopOfDescent, PIREPStatusLanding,
		PIREPStatusLanded, PIREPStatusArrived, PIREPStatusDeboarding,
		PIREPStatusPostShutdown:
		return true
	default:
		return false
	}
}
