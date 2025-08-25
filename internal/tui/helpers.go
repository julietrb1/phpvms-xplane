package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func conditionalDisplay(err *error) string {
	if *err != nil {
		return lipgloss.NewStyle().
			Foreground(colourAttention).
			Render((*err).Error())
	}
	return styleSecondary.
		Render("OK")
}

func conditionalAttentionString(input *string) string {
	if input == nil {
		return styleAttention.Render("(none)")
	}
	return *input
}

func conditionalAttentionTime(input *time.Time) string {
	if input == nil {
		return styleAttention.Render("(none)")
	}
	return input.Format(time.RFC3339)
}

func (model *Model) findSelectedAirline() string {
	var airlineInfo string
	for _, item := range model.airlineList.Items() {
		if airlineItem, ok := item.(AirlineItem); ok {
			if airlineItem.Airline.ID == model.selectedAirlineID {
				airlineInfo = fmt.Sprintf("%s (%s)",
					airlineItem.Airline.ICAO,
					airlineItem.Airline.Name)
				break
			}
		}
	}
	if airlineInfo == "" {
		airlineInfo = fmt.Sprintf("ID: %d", model.selectedAirlineID)
	}
	return airlineInfo
}
