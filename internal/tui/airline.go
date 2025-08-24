package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/julietrb1/phpvms-xplane/models"
	"io"
)

type AirlineItem struct {
	Airline models.Airline
}

func NewAirlineItem(airline models.Airline) AirlineItem {
	return AirlineItem{
		Airline: airline,
	}
}

func (i AirlineItem) Title() string {
	if i.Airline.ICAO != "" {
		return i.Airline.ICAO
	}
	return i.Airline.Name
}

func (i AirlineItem) Description() string {
	return i.Airline.Name
}

func (i AirlineItem) FilterValue() string {
	return i.Airline.FilterValue()
}

type AirlineDelegate struct{}

func (d AirlineDelegate) Height() int {
	return 2
}

func (d AirlineDelegate) Spacing() int {
	return 1
}

func (d AirlineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d AirlineDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(AirlineItem)
	if !ok {
		return
	}

	var title, desc string

	if m.Width() <= 0 {
		return
	}

	maxWidth := m.Width() - 4
	if maxWidth < 0 {
		maxWidth = 0
	}

	title = i.Title()
	if len(title) > maxWidth {
		title = title[:maxWidth-3] + "..."
	}

	desc = i.Description()
	if len(desc) > maxWidth {
		desc = desc[:maxWidth-3] + "..."
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	if index == m.Index() {
		title = selectedStyle.Render(title)
		desc = selectedStyle.Render(desc)
	} else {
		title = normalStyle.Render(title)
		desc = descStyle.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

func ConvertToAirlineItems(data []models.Airline) []list.Item {
	items := make([]list.Item, 0, len(data))
	for _, airline := range data {
		item := NewAirlineItem(airline)
		items = append(items, item)
	}
	return items
}

func GetSelectedAirlineID(model list.Model) int {
	if model.SelectedItem() == nil {
		return 0
	}

	item, ok := model.SelectedItem().(AirlineItem)
	if !ok {
		return 0
	}

	return item.Airline.ID
}
