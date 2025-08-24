package models

import "fmt"

type PirepState int

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
