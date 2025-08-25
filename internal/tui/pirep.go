package tui

import (
	"fmt"
	"math"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/julietrb1/phpvms-xplane/internal/api"
)

func (model *Model) startPIREP() tea.Cmd {
	return func() tea.Msg {
		if model.selectedAircraftID <= 0 {
			model.statusMessage = "Aircraft required"
			return nil
		}

		level, err := strconv.Atoi(model.flightInputs[6].Value())
		if err != nil {
			model.statusMessage = "Invalid altitude"
			return nil
		}

		plannedDistance, err := strconv.Atoi(model.flightInputs[5].Value())
		if err != nil {
			model.statusMessage = "Invalid planned distance"
			return nil
		}

		plannedFlightTime, err := strconv.Atoi(model.flightInputs[8].Value())
		if err != nil {
			model.statusMessage = "Invalid planned flight time"
			return nil
		}

		blockFuel, err := strconv.Atoi(model.flightInputs[7].Value())
		if err != nil {
			model.statusMessage = "Invalid block fuel"
			return nil
		}

		data := api.PrefilePIREPRequest{
			AirlineID:          model.selectedAirlineID,
			AircraftID:         model.selectedAircraftID,
			FlightType:         "J", // Scheduled Pax
			FlightNumber:       model.flightInputs[0].Value(),
			DepartureAirportID: model.flightInputs[1].Value(),
			ArrivalAirportID:   model.flightInputs[2].Value(),
			AlternateAirportID: model.flightInputs[3].Value(),
			Route:              model.flightInputs[9].Value(),
			Level:              level,
			PlannedDistance:    plannedDistance,
			PlannedFlightTime:  plannedFlightTime,
			BlockFuel:          blockFuel,
			Source:             1,
			SourceName:         "vmsacars",
			Fields: map[string]interface{}{
				"Simulator":              "X-Plane 12",
				"Unlimited Fuel":         "Off",
				"Network Online":         "VATSIM",
				"Network Callsign Check": "0",
				"Network Callsign Used":  model.flightInputs[4].Value(),
			},
		}

		pirepID, err := model.flightService.Prefile(model.ctx, data)
		return prefileDataMsg{pirepID, err}
	}
}

func (model *Model) filePIREP() tea.Cmd {
	return func() tea.Msg {
		snapshot := model.metrics.Snapshot()
		lastFuel := *snapshot.LastFuel
		if lastFuel == 0 {
			return pirepFiledMsg{
				error: fmt.Errorf("no last fuel level"),
			}
		}
		startingFuel, err := strconv.Atoi(model.flightInputs[7].Value())
		if err != nil {
			return pirepFiledMsg{error: fmt.Errorf("invalid starting fuel")}
		}
		fuelUsed := int(math.Max(0, float64(startingFuel-lastFuel)))
		if fuelUsed == 0 {
			return pirepFiledMsg{error: fmt.Errorf("no fuel used")}
		}
		data := api.FilePIREPRequest{
			FlightTime:  *snapshot.LastFlightTime,
			FuelUsedLbs: int(math.Ceil(float64(fuelUsed) * 2.20462)),
			Distance:    *snapshot.LastDistance,
		}
		err = model.flightService.FileFlight(model.ctx, data)
		return pirepFiledMsg{error: err}
	}
}

func (model *Model) cancelPIREP() tea.Cmd {
	return func() tea.Msg {
		return pirepCancelledMsg{error: model.flightService.CancelFlight(model.ctx)}
	}
}
