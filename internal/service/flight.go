package service

import (
	"context"
	"fmt"
	"github.com/julietrb1/phpvms-xplane/internal/api"
	"github.com/julietrb1/phpvms-xplane/internal/udp"
	"github.com/julietrb1/phpvms-xplane/models"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
)

type FlightService struct {
	Client        *api.Client
	Logger        *slog.Logger
	StateMachine  *StateMachine
	ActivePirepID atomic.Pointer[string]
	InitialFuel   int
	fuelMutex     sync.Mutex
}

func NewFlightService(client *api.Client, logger *slog.Logger) *FlightService {
	if logger == nil {
		logger = slog.Default()
	}

	s := &FlightService{
		Client:       client,
		Logger:       logger,
		StateMachine: NewStateMachine(),
	}
	s.ActivePirepID.Store(nil)
	return s
}

func (service *FlightService) Prefile(ctx context.Context, flightData api.PrefilePIREPRequest) (*string, error) {
	if id := service.ActivePirepID.Load(); id != nil && *id != "" {
		return nil, fmt.Errorf("there is already an active PIREP (ID: %s)", *id)
	}

	result, err := service.Client.PrefilePIREP(ctx, flightData)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, fmt.Errorf("failed to get PIREP ID from response")
	}

	service.SetActivePirepID(result.Data.ID)
	service.StateMachine.SetState(PIREPStateInProgress)

	return &result.Data.ID, nil
}

func (service *FlightService) UpdateFlight(ctx context.Context, status string, distance int, fuelRemainingKG int, flightTimeMin int) error {
	pirepID := service.ActivePirepID.Load()
	if pirepID == nil {
		service.Logger.Debug("No active PIREP, skipping update")
		return nil
	}

	if !service.StateMachine.CanUpdate() {
		service.Logger.Warn("PIREP is in read-only state, skipping update",
			"pirep_id", pirepID,
			"state", service.StateMachine.CurrentState.String())
		return nil
	}

	if status != "" && !ValidateStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	var data api.FlightUpdateRequest
	service.fuelMutex.Lock()
	initialFuel := service.InitialFuel
	if initialFuel > 0 {
		fuelUsed := int(math.Max(0, float64(initialFuel-fuelRemainingKG)))
		data = api.FlightUpdateRequest{
			Status:      status,
			Distance:    distance,
			FuelUsedLbs: int(math.Ceil(float64(fuelUsed) * 2.20462)),
			FlightTime:  flightTimeMin,
		}
	}
	if initialFuel == 0 {
		service.InitialFuel = fuelRemainingKG
	}
	service.fuelMutex.Unlock()

	if err := service.Client.UpdatePIREP(ctx, *pirepID, data); err != nil {
		return fmt.Errorf("failed to update PIREP: %w", err)
	}
	return nil
}

func (service *FlightService) SendPosition(ctx context.Context, pos udp.Position) error {
	pirepID := service.ActivePirepID.Load()
	if pirepID == nil {
		service.Logger.Debug("No active PIREP, skipping position update")
		return nil
	}

	if !service.StateMachine.CanUpdate() {
		service.Logger.Warn("PIREP is in read-only state, skipping position update",
			"pirep_id", pirepID,
			"state", service.StateMachine.CurrentState.String())
		return nil
	}

	if !isValidLatitude(pos.Lat) || !isValidLongitude(pos.Lon) {
		return fmt.Errorf("invalid position coordinates: lat=%f, lon=%f", pos.Lat, pos.Lon)
	}

	data := api.PositionUpdateRequest{
		Lat:        pos.Lat,
		Lon:        pos.Lon,
		AltMSL:     pos.AltMSL,
		AltAGL:     pos.AltAGL,
		GS:         pos.GS,
		SimTime:    pos.SimTime,
		DistanceNM: pos.DistanceNM,
		Heading:    pos.Heading,
		IAS:        pos.IAS,
		VSFPM:      pos.VSFPM,
	}

	if err := service.Client.PostACARSPosition(ctx, *pirepID, data); err != nil {
		return fmt.Errorf("failed to send ACARS position: %w", err)
	}
	return nil
}

func (service *FlightService) FileFlight(ctx context.Context, data api.FilePIREPRequest) error {
	pirepID := service.ActivePirepID.Load()
	if pirepID == nil {
		return fmt.Errorf("no active PIREP to file")
	}

	if !service.StateMachine.CanFile() {
		return fmt.Errorf("PIREP cannot be filed in current state: %s", service.StateMachine.CurrentState.String())
	}

	if err := service.Client.FilePIREP(ctx, *pirepID, data); err != nil {
		return fmt.Errorf("failed to file PIREP: %w", err)
	}

	service.StateMachine.SetState(PIREPStatePending)
	service.ResetActivePirep()

	return nil
}

func (service *FlightService) CancelFlight(ctx context.Context) error {
	pirepID := service.ActivePirepID.Load()
	if pirepID == nil {
		return fmt.Errorf("no active PIREP to cancel")
	}

	if !service.StateMachine.CanCancel() {
		return fmt.Errorf("PIREP cannot be cancelled in current state: %s", service.StateMachine.CurrentState.String())
	}

	if err := service.Client.CancelPIREP(ctx, *pirepID); err != nil {
		return fmt.Errorf("failed to cancel PIREP: %w", err)
	}

	service.StateMachine.SetState(PIREPStateCancelled)

	service.ActivePirepID.Store(nil)
	service.fuelMutex.Lock()
	service.InitialFuel = 0
	service.fuelMutex.Unlock()

	return nil
}

func (service *FlightService) SetActivePirepID(id string) {
	service.ActivePirepID.Store(&id)
}

func (service *FlightService) ResetActivePirep() {
	service.ActivePirepID.Store(nil)
	service.fuelMutex.Lock()
	service.InitialFuel = 0
	service.fuelMutex.Unlock()
	service.StateMachine.SetState(PIREPStateInProgress)
}

func (service *FlightService) GetAPIClient() *api.Client {
	return service.Client
}

func (service *FlightService) GetActivePirepID() *string {
	return service.ActivePirepID.Load()
}

func (service *FlightService) HandlePayload(ctx context.Context, payload *udp.Payload) (error, error) {
	pirepID := service.ActivePirepID.Load()
	if pirepID == nil {
		noPirepIDErr := fmt.Errorf("no active PIREP")
		return noPirepIDErr, noPirepIDErr
	}

	updateFlightsErr := service.UpdateFlight(ctx, payload.Status, payload.Position.DistanceNM, payload.Fuel, payload.FlightTime)
	updatePositionErr := service.SendPosition(ctx, payload.Position)

	return updateFlightsErr, updatePositionErr
}

func (service *FlightService) GetAirlines(ctx context.Context) ([]models.Airline, error) {
	airlines, err := service.Client.GetAirlines(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Implement pagination
	return airlines.Data, nil
}

func (service *FlightService) ListPIREPs(ctx context.Context) ([]models.ListedPIREP, error) {
	response, err := service.Client.ListPIREPs(ctx)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no active PIREPs found")
	}

	if response.Data[0].State != int(models.PIREPStateInProgress) {
		return nil, fmt.Errorf("latest PIREP not in progress")
	}

	// TODO: Implement pagination
	return response.Data, nil
}

func (service *FlightService) GetUserAircraftList(ctx context.Context) ([]models.Aircraft, error) {
	aircraftResponse, err := service.Client.GetUserFleet(ctx)
	if err != nil {
		return nil, err
	}

	if len(aircraftResponse.Data) == 0 {
		return nil, fmt.Errorf("no aircraft found")
	}

	// TODO: Implement pagination
	return aircraftResponse.Data[0].Aircraft, nil
}

func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func isValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}
