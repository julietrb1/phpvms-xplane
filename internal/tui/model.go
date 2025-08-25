package tui

import (
	"context"
	"fmt"
	"github.com/julietrb1/phpvms-xplane/internal/api"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/julietrb1/phpvms-xplane/internal/config"
	"github.com/julietrb1/phpvms-xplane/internal/service"
	"github.com/julietrb1/phpvms-xplane/internal/udp"
)

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 0, 1, 2).
			Foreground(colourPrimary)
	colourPrimary  = lipgloss.Color("205")
	styleSecondary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170"))
	colourBorder     = lipgloss.Color("240")
	colourSubtle     = lipgloss.Color("245")
	colourAttention  = lipgloss.Color("214")
	colourText       = lipgloss.Color("229")
	colourBackground = lipgloss.Color("57")
	headingStyle     = lipgloss.NewStyle().
				Bold(true).
				Underline(true)
	stylePairValue = lipgloss.NewStyle().Width(22).
			Foreground(colourSubtle)
	styleAttention = lipgloss.NewStyle().
			Foreground(colourAttention)
)

type keyMap struct {
	Help             key.Binding
	Quit             key.Binding
	Start            key.Binding
	File             key.Binding
	Cancel           key.Binding
	Reset            key.Binding
	Tab              key.Binding
	ShiftTab         key.Binding
	Enter            key.Binding
	Back             key.Binding
	SelectAircraft   key.Binding
	SelectAirline    key.Binding
	FetchSimbrief    key.Binding
	FetchActivePIREP key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Start, k.File, k.Cancel, k.Reset, k.SelectAircraft, k.SelectAirline, k.FetchSimbrief, k.FetchActivePIREP}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
		{k.Start, k.File, k.Cancel, k.Reset},
		{k.Enter, k.Back},
		{k.SelectAircraft, k.SelectAirline, k.FetchSimbrief, k.FetchActivePIREP},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Start: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prefile"),
	),
	File: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "file"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cancel active PIREP"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset active PIREP"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
	),
	SelectAircraft: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "select aircraft"),
	),
	SelectAirline: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "select airline"),
	),
	FetchSimbrief: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "fetch SimBrief OFP"),
	),
	FetchActivePIREP: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "fetch pending PIREP"),
	),
}

type Model struct {
	ctx                context.Context
	cancel             context.CancelFunc
	metrics            *udp.Metrics
	flightService      *service.FlightService
	logger             *slog.Logger
	help               help.Model
	spinner            spinner.Model
	keys               keyMap
	width              int
	height             int
	ready              bool
	showHelp           bool
	lastUpdate         time.Time
	statusMessage      string
	activeTab          int
	flightInputs       []textinput.Model
	flightTable        table.Model
	aircraftList       list.Model
	showAircraftList   bool
	selectedAircraftID int
	airlineList        list.Model
	showAirlineList    bool
	selectedAirlineID  int
	config             *config.Config
}

func NewModel(ctx context.Context, cancel context.CancelFunc, metrics *udp.Metrics, flightService *service.FlightService, cfg *config.Config, logger *slog.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colourPrimary)

	flightInputs := make([]textinput.Model, 10)
	for i := range flightInputs {
		t := textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Flight number"
			t.SetValue("1")
			t.Validate = func(s string) error {
				if number, err := strconv.Atoi(s); err != nil || number < 1 || number > 9999 {
					return fmt.Errorf("flight number must be a number between 1 and 9999")
				}
				return nil
			}
		case 1:
			t.Placeholder = "Departure"
			t.CharLimit = 4
		case 2:
			t.Placeholder = "Arrival"
			t.CharLimit = 4
		case 3:
			t.Placeholder = "Alternate"
			t.CharLimit = 4
		case 4:
			t.Placeholder = "e.g. QFA1"
			t.CharLimit = 6
		case 5:
			t.Placeholder = "e.g. 250"
			t.CharLimit = 4
		case 6:
			t.Placeholder = "e.g. 41000"
			t.CharLimit = 5
		case 7:
			t.Placeholder = "e.g. 1080"
			t.CharLimit = 6
		case 8:
			t.Placeholder = "e.g. 95"
			t.CharLimit = 3
		case 9:
			t.Placeholder = "e.g. DCT"
			t.CharLimit = 200
		}

		flightInputs[i] = t
	}

	columns := []table.Column{
		{Title: "Field", Width: 15},
		{Title: "Value", Width: 25},
	}
	rows := []table.Row{
		{"Flight Number", ""},
		{"Aircraft", ""},
		{"Departure", ""},
		{"Arrival", ""},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colourBorder).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(colourText).
		Background(colourBackground).
		Bold(false)
	t.SetStyles(tableStyles)

	aircraftDelegate := AircraftDelegate{}
	aircraftList := list.New([]list.Item{}, aircraftDelegate, 0, 0)
	aircraftList.Title = "Select Aircraft"
	aircraftList.SetShowStatusBar(false)
	aircraftList.SetFilteringEnabled(true)
	aircraftList.Styles.Title = styleTitle

	airlineDelegate := AirlineDelegate{}
	airlineList := list.New([]list.Item{}, airlineDelegate, 0, 0)
	airlineList.Title = "Select Airline"
	airlineList.SetShowStatusBar(false)
	airlineList.SetFilteringEnabled(true)
	airlineList.Styles.Title = styleTitle

	selectedAircraftID := 0
	selectedAirlineID := 0

	if cfg != nil {
		selectedAircraftID = cfg.SelectedAircraftID
		selectedAirlineID = cfg.SelectedAirlineID
	}

	return Model{
		ctx:                ctx,
		cancel:             cancel,
		metrics:            metrics,
		flightService:      flightService,
		logger:             logger,
		help:               help.New(),
		spinner:            s,
		keys:               keys,
		lastUpdate:         time.Now(),
		flightInputs:       flightInputs,
		flightTable:        t,
		aircraftList:       aircraftList,
		showAircraftList:   false,
		selectedAircraftID: selectedAircraftID,
		airlineList:        airlineList,
		showAirlineList:    false,
		selectedAirlineID:  selectedAirlineID,
		config:             cfg,
		statusMessage:      "Hi!",
	}
}

func (model *Model) Init() tea.Cmd {
	return tea.Batch(
		model.spinner.Tick,
		tickCmd(),
		model.fetchAircraftList(),
		model.fetchAirlineList(),
	)
}

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
		model.statusMessage = "Loading SimBrief OFP data..."
		apiClient := model.flightService.GetAPIClient()
		ofpData, err := apiClient.GetSimbriefOFP(model.ctx, model.config.SimbriefUserID)
		if err != nil {
			model.logger.Error("Failed to fetch SimBrief OFP data", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to fetch SimBrief OFP data: %v", err)
			return nil
		}

		routeDistance, err := strconv.Atoi(ofpData.General.RouteDistance)
		if err != nil {
			model.logger.Error("Failed to parse route distance", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to parse route distance: %v", err)
			return nil
		}

		initialAltitude, err := strconv.Atoi(ofpData.General.InitialAltitude)
		if err != nil {
			model.logger.Error("Failed to parse initial altitude", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to parse initial altitude: %v", err)
			return nil
		}

		blockFuel, err := strconv.Atoi(ofpData.Fuel.PlanRamp)
		if err != nil {
			model.logger.Error("Failed to parse block fuel", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to parse block fuel: %v", err)
			return nil
		}

		flightTime, err := strconv.Atoi(ofpData.Times.EstTimeEnroute)
		if err != nil {
			model.logger.Error("Failed to parse flight time", "error", err)
			model.statusMessage = fmt.Sprintf("Failed to parse flight time: %v", err)
			return nil
		}

		return simbriefDataMsg{
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

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if model.showAircraftList {
			switch {
			case key.Matches(msg, model.keys.Quit):
				model.cancel()
				return model, tea.Quit
			case key.Matches(msg, model.keys.Back):
				model.showAircraftList = false
				return model, nil
			case key.Matches(msg, model.keys.Enter):
				id := GetSelectedAircraftID(model.aircraftList)
				if id > 0 {
					model.selectedAircraftID = id
					model.showAircraftList = false
					model.statusMessage = fmt.Sprintf("Selected aircraft ID: %d", id)

					if model.config != nil {
						model.config.SelectedAircraftID = id
						if err := model.config.SavePreferences(""); err != nil {
							model.logger.Error("Failed to save preferences", "error", err)
						}
					}
				}
				return model, nil
			default:
				var cmd tea.Cmd
				model.aircraftList, cmd = model.aircraftList.Update(msg)
				return model, cmd
			}
		}

		if model.showAirlineList {
			switch {
			case key.Matches(msg, model.keys.Quit):
				model.cancel()
				return model, tea.Quit
			case key.Matches(msg, model.keys.Back):
				model.showAirlineList = false
				return model, nil
			case key.Matches(msg, model.keys.Enter):
				id := GetSelectedAirlineID(model.airlineList)
				if id > 0 {
					model.selectedAirlineID = id
					model.showAirlineList = false
					model.statusMessage = fmt.Sprintf("Selected airline ID: %d", id)

					if model.config != nil {
						model.config.SelectedAirlineID = id
						if err := model.config.SavePreferences(""); err != nil {
							model.logger.Error("Failed to save preferences", "error", err)
						}
					}
				}
				return model, nil
			default:
				var cmd tea.Cmd
				model.airlineList, cmd = model.airlineList.Update(msg)
				return model, cmd
			}
		}

		var focusedFlightInput *int
		if model.activeTab == 0 {
			for i := range model.flightInputs {
				if model.flightInputs[i].Focused() {
					focusedFlightInput = &i
					break
				}
			}

			if focusedFlightInput != nil && key.Matches(msg, model.keys.Back) {
				model.flightInputs[*focusedFlightInput].Blur()
				break
			} else if focusedFlightInput != nil && !key.Matches(msg, model.keys.Tab) && !key.Matches(msg, model.keys.ShiftTab) && !key.Matches(msg, model.keys.Quit) {
				break
			}
		}

		switch {
		case key.Matches(msg, model.keys.Quit):
			model.cancel()
			return model, tea.Quit
		case key.Matches(msg, model.keys.Help):
			model.showHelp = !model.showHelp
		case key.Matches(msg, model.keys.SelectAircraft):
			model.showAircraftList = true
			if len(model.aircraftList.Items()) == 0 {
				return model, model.fetchAircraftList()
			}
		case key.Matches(msg, model.keys.SelectAirline):
			model.showAirlineList = true
			if len(model.airlineList.Items()) == 0 {
				return model, model.fetchAirlineList()
			}
		case key.Matches(msg, model.keys.FetchSimbrief):
			if model.config.SimbriefUserID != "" {
				model.statusMessage = "Fetching SimBrief OFP..."
				return model, model.fetchSimbriefData()
			} else {
				model.statusMessage = "Set SIMBRIEF_USER_ID in .env"
			}
		case key.Matches(msg, model.keys.FetchActivePIREP):
			model.statusMessage = "Fetching active PIREP..."
			return model, model.fetchInProgressPIREP()
		case key.Matches(msg, model.keys.Start):
			if model.activeTab == 0 {
				model.statusMessage = "Prefiling PIREP..."
				cmd := model.startPIREP()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		case key.Matches(msg, model.keys.File):
			if model.activeTab == 0 {
				model.statusMessage = "Filing PIREP..."
				cmd := model.filePIREP()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		case key.Matches(msg, model.keys.Cancel):
			if model.activeTab == 0 {
				model.statusMessage = "Cancelling PIREP..."
				cmd := model.cancelPIREP()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		case key.Matches(msg, model.keys.Reset):
			if model.activeTab == 0 {
				model.flightService.ResetActivePirep()
				model.statusMessage = "Active PIREP reset"
			}
		case key.Matches(msg, model.keys.Tab):
			if model.activeTab == 0 {
				if focusedFlightInput != nil {
					model.flightInputs[*focusedFlightInput].Blur()
					if *focusedFlightInput < len(model.flightInputs)-1 {
						model.flightInputs[(*focusedFlightInput + 1)].Focus()
					}
				} else {
					model.flightInputs[0].Focus()
				}
			}
		case key.Matches(msg, model.keys.ShiftTab):
			if model.activeTab == 0 {
				if focusedFlightInput != nil {
					model.flightInputs[*focusedFlightInput].Blur()
					if *focusedFlightInput > 0 {
						model.flightInputs[(*focusedFlightInput - 1)].Focus()
					}
				} else {
					model.flightInputs[len(model.flightInputs)-1].Focus()
				}
			}
		}

	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		model.ready = true
		model.help.Width = msg.Width

		top, right, bottom, left := 2, 2, 2, 2
		model.aircraftList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		model.airlineList.SetSize(msg.Width-left-right, msg.Height-top-bottom)

	case tickMsg:
		model.lastUpdate = time.Time(msg)
		cmds = append(cmds, tickCmd())

	case spinner.TickMsg:
		var cmd tea.Cmd
		model.spinner, cmd = model.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case aircraftListUpdatedMsg:
		model.aircraftList.SetItems(msg.items)
		if len(msg.items) > 0 {
			model.statusMessage = fmt.Sprintf("Loaded %d aircraft", len(msg.items))
		} else {
			model.statusMessage = "No aircraft found"
		}

	case airlineListUpdatedMsg:
		model.airlineList.SetItems(msg.items)
		if len(msg.items) > 0 {
			model.statusMessage = fmt.Sprintf("Loaded %d airlines", len(msg.items))
		} else {
			model.statusMessage = "No airlines found"
		}

	case simbriefDataMsg:
		if msg.origin != "" && msg.destination != "" {
			model.flightInputs[1].SetValue(msg.origin)
			model.flightInputs[2].SetValue(msg.destination)
			model.flightInputs[3].SetValue(msg.alternate)
			model.flightInputs[4].SetValue(msg.flightNumber)
			model.flightInputs[5].SetValue(strconv.Itoa(msg.planDist))
			model.flightInputs[6].SetValue(strconv.Itoa(msg.initialAltitude))
			model.flightInputs[7].SetValue(strconv.Itoa(msg.blockFuel))
			model.flightInputs[8].SetValue(strconv.Itoa(msg.flightTime))
			model.flightInputs[9].SetValue(msg.route)

			model.statusMessage = fmt.Sprintf("SimBrief OFP loaded: %s to %s", msg.origin, msg.destination)
		} else {
			model.statusMessage = "Failed to extract origin, destination, alternate from SimBrief OFP"
		}

	case fetchInProgressPIREPMsg:
		if msg.error != nil {
			model.statusMessage = fmt.Sprintf("Failed to fetch active PIREP: %v", msg.error)
		} else {
			model.flightService.SetActivePirepID(msg.pirep.ID)
			model.flightService.StateMachine.SetState(service.PIREPStateInProgress)
			model.flightInputs[1].SetValue(msg.pirep.DptAirportID)
			model.flightInputs[2].SetValue(msg.pirep.ArrAirportID)
			if msg.pirep.AltAirportID != nil {
				model.flightInputs[3].SetValue(*msg.pirep.AltAirportID)
			} else {
				model.flightInputs[3].SetValue("")
			}

			// TODO: Fix this
			//var fields map[string]string
			//if err := json.Unmarshal(msg.pirep.Fields, &fields); err != nil {
			//	model.statusMessage = fmt.Sprintf("Failed to parse PIREP fields: %v", err)
			//	break
			//}
			//if networkCallsign, exists := fields["Network Callsign Used"]; exists {
			//	model.flightInputs[4].SetValue(networkCallsign)
			//} else {
			//	model.flightInputs[4].SetValue("")
			//}

			model.flightInputs[5].SetValue(strconv.Itoa(int(msg.pirep.Distance.Nmi)))
			model.flightInputs[6].SetValue(strconv.Itoa(*msg.pirep.Level))
			model.flightInputs[8].SetValue(strconv.Itoa(msg.pirep.FlightTime))
			model.flightInputs[9].SetValue(msg.pirep.Route)

			model.statusMessage = "Active PIREP fetched"
		}

	case pirepCancelledMsg:
		if msg.error != nil {
			model.statusMessage = fmt.Sprintf("Failed to cancel PIREP: %v", msg.error)
		} else {
			model.statusMessage = "PIREP cancelled"
		}

	case pirepFiledMsg:
		if msg.error != nil {
			model.statusMessage = fmt.Sprintf("Failed to file PIREP: %v", msg.error)
		} else {
			model.statusMessage = "PIREP filed"
		}

	case prefileDataMsg:
		if msg.error != nil {
			model.statusMessage = fmt.Sprintf("Failed to prefile PIREP: %v", msg.error)
		} else {
			model.statusMessage = "PIREP prefiled"
		}
	}

	if model.activeTab == 0 && !model.showAircraftList && !model.showAirlineList {
		for i := range model.flightInputs {
			var cmd tea.Cmd
			model.flightInputs[i], cmd = model.flightInputs[i].Update(msg)
			if i > 0 && i < 5 {
				model.flightInputs[i].SetValue(strings.ToUpper(model.flightInputs[i].Value()))
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return model, tea.Batch(cmds...)
}

func conditionalDisplay(err *error) string {
	if *err != nil {
		return lipgloss.NewStyle().
			Foreground(colourAttention).
			Render((*err).Error())
	}
	return styleSecondary.
		Render("OK")
}

func (model *Model) View() string {
	if !model.ready {
		return "Spinning prop..."
	}

	if model.showAircraftList {
		return model.aircraftList.View()
	}
	if model.showAirlineList {
		return model.airlineList.View()
	}

	snapshot := model.metrics.Snapshot()

	var s string

	s += styleTitle.
		Render("PXP: the phpVMS ACARS Client") + "\n"

	s += styleSecondary.
		Render(model.statusMessage) + "\n\n"

	// ACARS transmissions

	s += headingStyle.
		Render("ACARS transmissions") + "\n"
	s += stylePairValue.
		Render("Last flight update:")
	s += conditionalDisplay(snapshot.UpdateFlightErr) + "\n"
	s += stylePairValue.
		Render("Last position update:")
	s += conditionalDisplay(snapshot.UpdatePositionErr) + "\n\n"

	s = model.renderUDPMetrics(s, snapshot)

	// Flight metrics

	s = model.renderFlightMetrics(s, snapshot)

	// Flight controls

	s += headingStyle.
		Render("Flight controls") + "\n"

	if model.activeTab == 0 {
		airlineLabel := stylePairValue.Render("Airline")
		if model.selectedAirlineID > 0 {
			airlineInfo := model.findSelectedAirline()
			s += fmt.Sprintf("%s %s\n", airlineLabel,
				styleSecondary.Render(airlineInfo))
		} else {
			s += fmt.Sprintf("%s %s\n", airlineLabel,
				lipgloss.NewStyle().Foreground(colourBorder).Render("Press 'l' to select airline"))
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
				stylePairValue.Render(label),
				input.View())
		}

		aircraftLabel := stylePairValue.Render("Aircraft")
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
			s += fmt.Sprintf("%s %s\n", aircraftLabel,
				styleSecondary.Render(aircraftInfo))
		} else {
			s += fmt.Sprintf("%s %s\n", aircraftLabel,
				lipgloss.NewStyle().Foreground(colourBorder).Render("Press 'a' to select aircraft"))
		}
	}

	helpView := model.help.View(model.keys)
	if model.showHelp {
		s += "\n" + helpView
	} else {
		s += "\n" + lipgloss.NewStyle().Render("Press ? for help")
	}

	return s
}

func conditionalAttentionString(input *string) string {
	if input == nil {
		return styleAttention.Render("(none)")
	}
	return styleAttention.Render(*input)
}

func conditionalAttentionTime(input *time.Time) string {
	if input == nil {
		return styleAttention.Render("(none)")
	}
	return input.Format(time.RFC3339)
}

func (model *Model) renderFlightMetrics(s string, snapshot udp.MetricsSnapshot) string {
	s += headingStyle.Render("Flight metrics") + "\n"

	if snapshot.LastStatus != nil {
		s += fmt.Sprintf("Last status: %s\n", *snapshot.LastStatus)
	}
	s += stylePairValue.Render("Fuel:")
	s += fmt.Sprintf("%d kg\n", *snapshot.LastFuel)

	s += stylePairValue.Render("Flight time:")
	s += fmt.Sprintf("%d minutes\n", *snapshot.LastFlightTime)

	s += stylePairValue.Render("Distance:")
	s += fmt.Sprintf("%d nm\n", *snapshot.LastDistance)

	pirepID := model.flightService.GetActivePirepID()
	s += stylePairValue.Render("Active PIREP ID:")
	s += fmt.Sprintf("%s\n\n", conditionalAttentionString(pirepID))
	return s
}

func (model *Model) renderUDPMetrics(s string, snapshot udp.MetricsSnapshot) string {
	s += headingStyle.
		Render("UDP metrics") + "\n"

	s += stylePairValue.Render("Packets:")
	if snapshot.PacketsAny == 0 {
		s += styleAttention.Render("(none)") + "\n"
	} else {
		s += fmt.Sprintf("%d total, %d errors",
			snapshot.PacketsAny, snapshot.PacketsErr) + "\n"
	}

	s += stylePairValue.Render("Last sender:")
	s += conditionalAttentionString(snapshot.LastSender) + "\n"

	s += stylePairValue.Render("Last packet:")
	s += conditionalAttentionTime(snapshot.LastPacketTime) + "\n\n"
	return s
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
				fmt.Errorf("no last fuel level"),
			}
		}
		startingFuel, err := strconv.Atoi(model.flightInputs[7].Value())
		if err != nil {
			return pirepFiledMsg{fmt.Errorf("invalid starting fuel")}
		}
		fuelUsed := int(math.Max(0, float64(startingFuel-lastFuel)))
		if fuelUsed == 0 {
			return pirepFiledMsg{fmt.Errorf("no fuel used")}
		}
		data := api.FilePIREPRequest{
			FlightTime:  *snapshot.LastFlightTime,
			FuelUsedLbs: int(math.Ceil(float64(fuelUsed) * 2.20462)),
			Distance:    *snapshot.LastDistance,
		}
		err = model.flightService.FileFlight(model.ctx, data)
		return pirepFiledMsg{err}
	}
}

func (model *Model) cancelPIREP() tea.Cmd {
	return func() tea.Msg {
		return pirepCancelledMsg{model.flightService.CancelFlight(model.ctx)}
	}
}
