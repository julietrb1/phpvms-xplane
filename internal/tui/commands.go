package tui

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) fetchAirlineList() tea.Cmd {
	return func() tea.Msg {
		airlines, err := model.flightService.GetAirlines(model.ctx)
		if err != nil {
			model.logger.Error("Failed to fetch airline list", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to fetch airline list: %v", err)
			return nil
		}

		if len(airlines) == 0 {
			model.statusMessage = "No airlines found in the API response"
			return nil
		}

		items := ConvertToAirlineItems(airlines)

		if len(items) == 0 {
			model.statusMessage = "Failed to convert airlines data to list items"
			return nil
		}

		return airlineListUpdatedMsg{items: items}
	}
}

func (model *Model) fetchInProgressPIREP() tea.Cmd {
	return func() tea.Msg {
		pireps, err := model.flightService.ListPIREPs(model.ctx)
		if err != nil {
			return fetchInProgressPIREPMsg{
				error: err,
			}
		}
		return fetchInProgressPIREPMsg{
			pirep: pireps[0],
		}
	}
}

func (model *Model) fetchSimbriefData() tea.Cmd {
	return func() tea.Msg {
		apiClient := model.flightService.GetAPIClient()
		ofpData, err := apiClient.GetSimbriefOFP(model.ctx, model.config.SimbriefUserID)
		if err != nil {
			return fetchSimbriefOFPErrorMsg{
				err: fmt.Errorf("failed to fetch SimBrief OFP: %w", err),
			}
		}

		routeDistance, err := strconv.Atoi(ofpData.General.RouteDistance)
		if err != nil {
			return fetchSimbriefOFPErrorMsg{
				err: fmt.Errorf("failed to parse route distance: %w", err),
			}
		}

		initialAltitude, err := strconv.Atoi(ofpData.General.InitialAltitude)
		if err != nil {
			return fetchSimbriefOFPErrorMsg{
				err: fmt.Errorf("failed to parse initial altitude: %w", err),
			}
		}

		blockFuel, err := strconv.Atoi(ofpData.Fuel.PlanRamp)
		if err != nil {
			return fetchSimbriefOFPErrorMsg{
				err: fmt.Errorf("failed to parse block fuel: %w", err),
			}
		}

		flightTime, err := strconv.Atoi(ofpData.Times.EstTimeEnroute)
		if err != nil {
			return fetchSimbriefOFPErrorMsg{
				err: fmt.Errorf("failed to parse flight time: %w", err),
			}
		}

		return fetchSimbriefOFPMsg{
			origin:          string(ofpData.Origin.ICAOCode),
			destination:     string(ofpData.Destination.ICAOCode),
			alternate:       string(ofpData.Alternate.ICAOCode),
			flightNumber:    ofpData.General.FlightNumber,
			planDist:        routeDistance,
			initialAltitude: initialAltitude,
			blockFuel:       blockFuel,
			flightTime:      flightTime,
			route:           ofpData.General.Route,
		}
	}
}

func (model *Model) fetchAircraftList() tea.Cmd {
	return func() tea.Msg {
		fleet, err := model.flightService.GetUserAircraftList(model.ctx)
		if err != nil {
			model.logger.Error("Failed to fetch aircraft list", "error", err)
			return nil
		}
		items := ConvertToAircraftItems(fleet)
		return aircraftListUpdatedMsg{items: items}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
