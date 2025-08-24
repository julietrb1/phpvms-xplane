package udp

type Payload struct {
	Status     string   `json:"status"`
	Position   Position `json:"position"`
	Fuel       int      `json:"fuel"`        // kg remaining
	FlightTime int      `json:"flight_time"` // minutes
}

type Position struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	AltMSL     int     `json:"altitude_msl"`
	AltAGL     int     `json:"altitude_agl"`
	GS         int     `json:"gs"`
	SimTime    string  `json:"sim_time"`
	DistanceNM int     `json:"distance"`
	Heading    int     `json:"heading"`
	IAS        int     `json:"ias"`
	VSFPM      int     `json:"vs"`
}
