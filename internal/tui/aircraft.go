package tui

import (
	"fmt"
	"github.com/julietrb1/phpvms-xplane/models"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AircraftItem struct {
	Aircraft models.Aircraft
}

func NewAircraftItem(aircraft models.Aircraft) AircraftItem {
	return AircraftItem{
		Aircraft: aircraft,
	}
}

func (i AircraftItem) Title() string {
	return i.Aircraft.Registration
}

func (i AircraftItem) Description() string {
	return fmt.Sprintf("%s - %s", i.Aircraft.ICAO, i.Aircraft.Name)
}

func (i AircraftItem) FilterValue() string {
	return i.Aircraft.FilterValue()
}

type AircraftDelegate struct{}

func (d AircraftDelegate) Height() int {
	return 2
}

func (d AircraftDelegate) Spacing() int {
	return 1
}

func (d AircraftDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d AircraftDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(AircraftItem)
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

func ConvertToAircraftItems(data []models.Aircraft) []list.Item {
	items := make([]list.Item, len(data))
	if len(data) == 0 {
		return items
	}
	for i, ac := range data {
		items[i] = NewAircraftItem(ac)
	}
	sortAircraftItems(items)
	return items
}

func sortAircraftItems(items []list.Item) {
	sort.Slice(items, func(i, j int) bool {
		regI := items[i].(AircraftItem).Aircraft.Registration
		regJ := items[j].(AircraftItem).Aircraft.Registration

		return strings.Compare(regI, regJ) < 0
	})
}

func GetSelectedAircraftID(model list.Model) int {
	if model.SelectedItem() == nil {
		return 0
	}
	item, ok := model.SelectedItem().(AircraftItem)
	if !ok {
		return 0
	}
	return item.Aircraft.ID
}
