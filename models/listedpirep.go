package models

import (
	"encoding/json"
	"time"
)

type Distances struct {
	M   float64 `json:"m"`
	Km  float64 `json:"km"`
	Mi  float64 `json:"mi"`
	Nmi float64 `json:"nmi"`
}

type Weights struct {
	Kg  float64 `json:"kg"`
	Lbs float64 `json:"lbs"`
}

type Airport struct {
	ID                 string     `json:"id"`
	IATA               string     `json:"iata"`
	ICAO               string     `json:"icao"`
	Name               string     `json:"name"`
	Location           string     `json:"location"`
	Region             string     `json:"region"`
	Country            string     `json:"country"`
	Timezone           string     `json:"timezone"`
	Hub                bool       `json:"hub"`
	Notes              *string    `json:"notes"`
	Latitude           float64    `json:"lat"`
	Longitude          float64    `json:"lon"`
	Elevation          string     `json:"elevation"`
	GroundHandlingCost *int       `json:"ground_handling_cost"`
	Fuel100LLCost      *int       `json:"fuel_100ll_cost"`
	FuelJetACost       *int       `json:"fuel_jeta_cost"`
	FuelMOGASCost      *int       `json:"fuel_mogas_cost"`
	DeletedAt          *time.Time `json:"deleted_at"`
	Description        string     `json:"description"`
}

type ListedPIREP struct {
	ID                string     `json:"id"`
	UserID            int        `json:"user_id"`
	AirlineID         int        `json:"airline_id"`
	AircraftID        int        `json:"aircraft_id"`
	EventID           *string    `json:"event_id"`
	FlightID          *string    `json:"flight_id"`
	FlightNumber      string     `json:"flight_number"`
	RouteCode         *string    `json:"route_code"`
	RouteLeg          *string    `json:"route_leg"`
	FlightType        string     `json:"flight_type"`
	DptAirportID      string     `json:"dpt_airport_id"`
	ArrAirportID      string     `json:"arr_airport_id"`
	AltAirportID      *string    `json:"alt_airport_id"`
	Level             *int       `json:"level"`
	Distance          Distances  `json:"distance"`
	PlannedDistance   Distances  `json:"planned_distance"`
	FlightTime        int        `json:"flight_time"`
	PlannedFlightTime *int       `json:"planned_flight_time"`
	ZFW               *int       `json:"zfw"`
	BlockFuel         Weights    `json:"block_fuel"`
	FuelUsed          Weights    `json:"fuel_used"`
	LandingRate       *float64   `json:"landing_rate"`
	Score             *int       `json:"score"`
	Route             string     `json:"route"`
	Notes             *string    `json:"notes"`
	Source            int        `json:"source"`
	SourceName        *string    `json:"source_name"`
	State             int        `json:"state"`
	Status            string     `json:"status"`
	SubmittedAt       time.Time  `json:"submitted_at"`
	BlockOffTime      time.Time  `json:"block_off_time"`
	BlockOnTime       time.Time  `json:"block_on_time"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	Airline           struct {
		ID      int    `json:"id"`
		ICAO    string `json:"icao"`
		IATA    string `json:"iata"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Logo    string `json:"logo"`
	} `json:"airline"`
	DptAirport Airport `json:"dpt_airport"`
	ArrAirport Airport `json:"arr_airport"`
	Ident      string  `json:"ident"`
	Phase      string  `json:"phase"`
	StatusText string  `json:"status_text"`
	Aircraft   struct {
		ID           int        `json:"id"`
		SubfleetID   int        `json:"subfleet_id"`
		Icao         string     `json:"icao"`
		Iata         string     `json:"iata"`
		AirportID    string     `json:"airport_id"`
		HubId        string     `json:"hub_id"`
		LandingTime  time.Time  `json:"landing_time"`
		Name         string     `json:"name"`
		Registration string     `json:"registration"`
		Fin          *string    `json:"fin"`
		HexCode      string     `json:"hex_code"`
		Selcal       *string    `json:"selcal"`
		DOW          Weights    `json:"dow"`
		MTOW         Weights    `json:"mtow"`
		MLW          Weights    `json:"mlw"`
		ZFW          Weights    `json:"zfw"`
		SimbriefType *string    `json:"simbrief_type"`
		FuelOnboard  Weights    `json:"fuel_onboard"`
		FlightTime   int        `json:"flight_time"`
		Status       *string    `json:"status"`
		State        int        `json:"state"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `json:"deleted_at"`
		Ident        string     `json:"ident"`
	} `json:"aircraft"`
	Fields json.RawMessage `json:"fields"`
}
