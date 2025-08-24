package models

import "time"

type AircraftFleet struct {
	Id                       int           `json:"id"`
	AirlineId                int           `json:"airline_id"`
	HubId                    *string       `json:"hub_id"`
	Type                     string        `json:"type"`
	SimbriefType             *string       `json:"simbrief_type"`
	Name                     string        `json:"name"`
	CostBlockHour            *string       `json:"cost_block_hour"`
	CostDelayMinute          *string       `json:"cost_delay_minute"`
	FuelType                 int           `json:"fuel_type"`
	GroundHandlingMultiplier int           `json:"ground_handling_multiplier"`
	CargoCapacity            *string       `json:"cargo_capacity"`
	FuelCapacity             *string       `json:"fuel_capacity"`
	GrossWeight              *string       `json:"gross_weight"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
	DeletedAt                *time.Time    `json:"deleted_at"`
	Fares                    []interface{} `json:"fares"`
	Aircraft                 []Aircraft    `json:"aircraft"`
}
