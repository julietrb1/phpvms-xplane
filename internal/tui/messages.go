package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/julietrb1/phpvms-xplane/models"
	"time"
)

type tickMsg time.Time

type aircraftListUpdatedMsg struct {
	items []list.Item
}

type airlineListUpdatedMsg struct {
	items []list.Item
}

type selectAircraftMsg struct {
	id int
}

type selectAirlineMsg struct {
	id int
}

type simbriefDataMsg struct {
	origin          string
	destination     string
	alternate       string
	flightNumber    string
	planDist        int
	initialAltitude int
	blockFuel       int
	flightTime      int
	route           string
}

type prefileDataMsg struct {
	pirepID *string
	error   error
}

type pirepCancelledMsg struct {
	error error
}

type pirepFiledMsg struct {
	error error
}

type fetchInProgressPIREPMsg struct {
	error error
	pirep models.ListedPIREP
}
