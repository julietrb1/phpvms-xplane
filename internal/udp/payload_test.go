package udp

import (
	"encoding/json"
	"testing"
)

func TestUDPPayloadDecode(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
		validate func(t *testing.T, payload *Payload)
	}{
		{
			name:     "minimal payload",
			jsonData: `{"status":"ENR"}`,
			wantErr:  false,
			validate: func(t *testing.T, payload *Payload) {
				if payload.Status != "ENR" {
					t.Errorf("Expected status ENR, got %s", payload.Status)
				}
				if payload.Position != nil {
					t.Errorf("Expected nil position, got %+v", payload.Position)
				}
				if payload.Fuel != nil {
					t.Errorf("Expected nil fuel, got %f", *payload.Fuel)
				}
				if payload.FlightTime != nil {
					t.Errorf("Expected nil flight time, got %f", *payload.FlightTime)
				}
				if len(payload.Events) != 0 {
					t.Errorf("Expected empty events, got %d events", len(payload.Events))
				}
			},
		},
		{
			name: "full payload",
			jsonData: `{
				"status": "ENR",
				"position": {
					"lat": 40.6398,
					"lon": -73.7789,
					"altitude_msl": 12000.0,
					"altitude_agl": 500.0,
					"gs": 320.0,
					"sim_time": 1724167500,
					"distance": 217.4,
					"heading": 255.0,
					"ias": 250.0,
					"vs": -500.0
				},
				"fuel": 4520.0,
				"flight_time": 83.0,
				"events": [{"log":"Passing 10,000 ft","sim_time":1724167400}]
			}`,
			wantErr: false,
			validate: func(t *testing.T, payload *Payload) {
				if payload.Status != "ENR" {
					t.Errorf("Expected status ENR, got %s", payload.Status)
				}

				if payload.Position == nil {
					t.Fatalf("Expected non-nil position")
				}
				if payload.Position.Lat != 40.6398 {
					t.Errorf("Expected lat 40.6398, got %f", payload.Position.Lat)
				}
				if payload.Position.Lon != -73.7789 {
					t.Errorf("Expected lon -73.7789, got %f", payload.Position.Lon)
				}

				if payload.Position.AltMSL == nil {
					t.Fatalf("Expected non-nil altitude_msl")
				}
				if *payload.Position.AltMSL != 12000.0 {
					t.Errorf("Expected altitude_msl 12000.0, got %f", *payload.Position.AltMSL)
				}

				if payload.Fuel == nil {
					t.Fatalf("Expected non-nil fuel")
				}
				if *payload.Fuel != 4520.0 {
					t.Errorf("Expected fuel 4520.0, got %f", *payload.Fuel)
				}

				if payload.FlightTime == nil {
					t.Fatalf("Expected non-nil flight_time")
				}
				if *payload.FlightTime != 83.0 {
					t.Errorf("Expected flight_time 83.0, got %f", *payload.FlightTime)
				}

				if len(payload.Events) != 1 {
					t.Fatalf("Expected 1 event, got %d", len(payload.Events))
				}
				if payload.Events[0].Log != "Passing 10,000 ft" {
					t.Errorf("Expected event log 'Passing 10,000 ft', got '%s'", payload.Events[0].Log)
				}
				if payload.Events[0].SimTime == nil {
					t.Fatalf("Expected non-nil event sim_time")
				}
				if *payload.Events[0].SimTime != 1724167400 {
					t.Errorf("Expected event sim_time 1724167400, got %d", *payload.Events[0].SimTime)
				}
			},
		},
		{
			name:     "scratch file example",
			jsonData: `{"fuel":694.7456831336,"status":"INI","flight_time":114.86914876302,"position":{"lat":-34.250967870196,"lon":148.24648374293,"altitude_msl":1272.1543825923,"altitude_agl":-0.0077407819366455,"gs":0.00066591917010024,"sim_time":1755770053,"distance":279.22966725162,"heading":20.255651473999,"vs":-0.023739677756384,"ias":9.6391868591309}}`,
			wantErr:  false,
			validate: func(t *testing.T, payload *Payload) {
				if payload.Status != "INI" {
					t.Errorf("Expected status INI, got %s", payload.Status)
				}

				if payload.Position == nil {
					t.Fatalf("Expected non-nil position")
				}
				if payload.Position.Lat != -34.250967870196 {
					t.Errorf("Expected lat -34.250967870196, got %f", payload.Position.Lat)
				}
				if payload.Position.Lon != 148.24648374293 {
					t.Errorf("Expected lon 148.24648374293, got %f", payload.Position.Lon)
				}

				if payload.Fuel == nil {
					t.Fatalf("Expected non-nil fuel")
				}
				if *payload.Fuel != 694.7456831336 {
					t.Errorf("Expected fuel 694.7456831336, got %f", *payload.Fuel)
				}

				if payload.FlightTime == nil {
					t.Fatalf("Expected non-nil flight_time")
				}
				if *payload.FlightTime != 114.86914876302 {
					t.Errorf("Expected flight_time 114.86914876302, got %f", *payload.FlightTime)
				}
			},
		},
		{
			name:     "invalid json",
			jsonData: `{"status":"ENR"`,
			wantErr:  true,
			validate: func(t *testing.T, payload *Payload) {
				// No validation needed for error case
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload Payload
			err := json.Unmarshal([]byte(tt.jsonData), &payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.validate(t, &payload)
			}
		})
	}
}
