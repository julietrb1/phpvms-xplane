package service

import "fmt"

type PirepState int
type PirepStatus string

const (
	PIREPStateInProgress PirepState = 0
	PIREPStatePending    PirepState = 1
	PIREPStateAccepted   PirepState = 2
	PIREPStateCancelled  PirepState = 3
	PIREPStateDeleted    PirepState = 4
	PIREPStateDraft      PirepState = 5
	PIREPStateRejected   PirepState = 6
	PIREPStatePaused     PirepState = 7
)

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

func (s PirepState) String() string {
	switch s {
	case PIREPStateInProgress:
		return "IN_PROGRESS"
	case PIREPStatePending:
		return "PENDING"
	case PIREPStateAccepted:
		return "ACCEPTED"
	case PIREPStateCancelled:
		return "CANCELLED"
	case PIREPStateDeleted:
		return "DELETED"
	case PIREPStateDraft:
		return "DRAFT"
	case PIREPStateRejected:
		return "REJECTED"
	case PIREPStatePaused:
		return "PAUSED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

func (s PirepState) IsReadOnly() bool {
	return s == PIREPStateAccepted || s == PIREPStateRejected || s == PIREPStateCancelled || s == PIREPStateDeleted
}

func (s PirepState) CanUpdate() bool {
	return !s.IsReadOnly()
}

func (s PirepState) CanCancel() bool {
	return s != PIREPStateAccepted && s != PIREPStateRejected && s != PIREPStateCancelled && s != PIREPStateDeleted
}

func (s PirepState) CanFile() bool {
	return s == PIREPStateInProgress || s == PIREPStateDraft || s == PIREPStatePaused
}

type StateMachine struct {
	CurrentState PirepState
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		CurrentState: PIREPStateInProgress,
	}
}

func (sm *StateMachine) SetState(state PirepState) {
	sm.CurrentState = state
}

func (sm *StateMachine) CanUpdate() bool {
	return sm.CurrentState.CanUpdate()
}

func (sm *StateMachine) CanCancel() bool {
	return sm.CurrentState.CanCancel()
}

func (sm *StateMachine) CanFile() bool {
	return sm.CurrentState.CanFile()
}

func (sm *StateMachine) IsReadOnly() bool {
	return sm.CurrentState.IsReadOnly()
}

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
