package tui

import (
	"fmt"

	"github.com/julietrb1/phpvms-xplane/internal/udp"
)

func (model *Model) renderFlightControls(s string) string {
	s += styleHeading.
		Render("Flight controls") + "\n"

	if model.activeTab == 0 {
		s += stylePairKey.Render("Airline")
		if model.selectedAirlineID > 0 {
			airlineInfo := model.findSelectedAirline()
			s += airlineInfo + "\n"
		} else {
			s += styleSecondary.Render("Press 'l' to select airline") + "\n"
		}

		for i, input := range model.flightInputs {
			var label string
			switch i {
			case 0:
				label = "ACARS flight no."
			case 1:
				label = "Departure"
			case 2:
				label = "Arrival"
			case 3:
				label = "Alternate"
			case 4:
				label = "SimBrief flight no."
			case 5:
				{
					label = "Plan dist. (nm)"
				}
			case 6:
				{
					label = "Altitude (ft)"
				}
			case 7:
				{
					label = "Block fuel (kg)"
				}
			case 8:
				{
					label = "Flight time (min)"
				}
			case 9:
				{
					label = "Route"
				}
			}
			s += fmt.Sprintf("%s %s\n",
				stylePairKey.Render(label),
				input.View())
		}

		s += stylePairKey.Render("Aircraft")
		if model.selectedAircraftID > 0 {
			var aircraftInfo string
			for _, item := range model.aircraftList.Items() {
				if aircraftItem, ok := item.(AircraftItem); ok {
					if aircraftItem.Aircraft.ID == model.selectedAircraftID {
						aircraftInfo = fmt.Sprintf("%s (%s - %s)",
							aircraftItem.Aircraft.Registration,
							aircraftItem.Aircraft.ICAO,
							aircraftItem.Aircraft.Name)
						break
					}
				}
			}
			if aircraftInfo == "" {
				aircraftInfo = fmt.Sprintf("ID: %d", model.selectedAircraftID)
			}
			s += aircraftInfo + "\n"
		} else {
			s += styleSecondary.Render("Press 'a' to select aircraft") + "\n"
		}
	}
	return s
}

func (model *Model) renderTitle(s string) string {
	s += styleTitle.
		Render("PXP: the phpVMS ACARS Client") + "\n"

	s += styleSubtitle.
		Render(model.statusMessage)
	return s
}

func (model *Model) renderACARSTransmissions(s string, snapshot udp.MetricsSnapshot) string {
	s += styleHeading.
		Render("ACARS transmissions") + "\n"

	pirepID := model.flightService.GetActivePirepID()
	s += stylePairKey.Render("Active PIREP ID:")
	s += fmt.Sprintf("%s\n", conditionalAttentionString(pirepID))

	s += stylePairKey.
		Render("Last flight update:")
	s += conditionalDisplay(snapshot.UpdateFlightErr) + "\n"

	s += stylePairKey.
		Render("Last position update:")
	s += conditionalDisplay(snapshot.UpdatePositionErr) + "\n"

	return s
}

func (model *Model) renderFlightMetrics(s string, snapshot udp.MetricsSnapshot) string {
	s += styleHeading.Render("Flight metrics") + "\n"

	if snapshot.LastStatus != nil {
		s += fmt.Sprintf("Last status: %s\n", *snapshot.LastStatus)
	}
	s += stylePairKey.Render("Fuel:")
	s += fmt.Sprintf("%d kg\n", *snapshot.LastFuel)

	s += stylePairKey.Render("Flight time:")
	s += fmt.Sprintf("%d minutes\n", *snapshot.LastFlightTime)

	s += stylePairKey.Render("Distance:")
	s += fmt.Sprintf("%d nm\n", *snapshot.LastDistance)

	return s
}

func (model *Model) renderUDPMetrics(s string, snapshot udp.MetricsSnapshot) string {
	s += styleHeading.
		Render("UDP metrics") + "\n"

	s += stylePairKey.Render("Packets:")
	if snapshot.PacketsAny == 0 {
		s += styleAttention.Render("(none)") + "\n"
	} else {
		s += fmt.Sprintf("%d total, %d errors",
			snapshot.PacketsAny, snapshot.PacketsErr) + "\n"
	}

	s += stylePairKey.Render("Last sender:")
	s += conditionalAttentionString(snapshot.LastSender) + "\n"

	s += stylePairKey.Render("Last packet:")
	s += conditionalAttentionTime(snapshot.LastPacketTime) + "\n"
	return s
}
