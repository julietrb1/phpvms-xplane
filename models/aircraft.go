package models

import "fmt"

type Aircraft struct {
	ID           int
	Registration string
	ICAO         string
	Name         string
}

func (a Aircraft) FilterValue() string {
	return fmt.Sprintf("%s %s %s", a.Registration, a.ICAO, a.Name)
}
