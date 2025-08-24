package models

import "fmt"

type Airline struct {
	ID      int
	ICAO    string
	IATA    string
	Name    string
	Country string
	Logo    string
}

func NewAirline(id int, icao, iata, name, country, logo string) Airline {
	return Airline{
		ID:      id,
		ICAO:    icao,
		IATA:    iata,
		Name:    name,
		Country: country,
		Logo:    logo,
	}
}

func (a Airline) FilterValue() string {
	return fmt.Sprintf("%s %s %s", a.ICAO, a.IATA, a.Name)
}
