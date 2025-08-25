package udp

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type Metrics struct {
	PacketsAny        atomic.Int64
	PacketsErr        atomic.Int64
	LastSender        atomic.Pointer[net.UDPAddr]
	LastPacketTime    atomic.Int64
	LastNonJSONHead   atomic.Pointer[string]
	LastStatus        atomic.Pointer[string]
	LastPosition      atomic.Pointer[Position]
	LastFuel          atomic.Int32
	LastFlightTime    atomic.Int32
	LastDistance      atomic.Int32
	UpdateFlightErr   atomic.Pointer[error]
	UpdatePositionErr atomic.Pointer[error]
}

func NewMetrics() *Metrics {
	initialFlightError := fmt.Errorf("(none)")
	initialPositionError := fmt.Errorf("(none)")
	metrics := &Metrics{}
	metrics.UpdateFlightErr.Store(&initialFlightError)
	metrics.UpdatePositionErr.Store(&initialPositionError)
	return metrics
}

type MetricsSnapshot struct {
	PacketsAny        int64      `json:"packets_any"`
	PacketsErr        int64      `json:"packets_err"`
	LastSender        *string    `json:"last_sender,omitempty"`
	LastPacketTime    *time.Time `json:"last_packet_time,omitempty"`
	LastNonJSONHead   string     `json:"last_non_json_head,omitempty"`
	LastStatus        *string    `json:"last_status,omitempty"`
	LastPosition      *Position  `json:"last_position,omitempty"`
	LastFuel          *int       `json:"last_fuel,omitempty"`
	LastFlightTime    *int       `json:"last_flight_time,omitempty"`
	LastDistance      *int       `json:"last_distance,omitempty"`
	UpdateFlightErr   *error     `json:"update_flight_err,omitempty"`
	UpdatePositionErr *error     `json:"update_position_err,omitempty"`
}

func (metrics *Metrics) Snapshot() MetricsSnapshot {
	var lastSender *string
	if addr := metrics.LastSender.Load(); addr != nil {
		lastSenderStr := addr.String()
		lastSender = &lastSenderStr
	}

	var lastPacketTime *time.Time
	if ts := metrics.LastPacketTime.Load(); ts > 0 {
		lastPacketTimeValue := time.Unix(ts, 0)
		lastPacketTime = &lastPacketTimeValue
	}

	var lastNonJSONHead string
	if head := metrics.LastNonJSONHead.Load(); head != nil {
		lastNonJSONHead = *head
	}

	var lastStatus *string
	if status := metrics.LastStatus.Load(); status != nil {
		lastStatus = status
	}

	var lastFuel int
	if fuel := metrics.LastFuel.Load(); fuel > 0 {
		lastFuel = int(fuel)
	}

	var lastFlightTime int
	if flightTime := metrics.LastFlightTime.Load(); flightTime > 0 {
		lastFlightTime = int(flightTime)
	}

	var lastDistance int
	if distance := metrics.LastDistance.Load(); distance > 0 {
		lastDistance = int(distance)
	}

	return MetricsSnapshot{
		PacketsAny:        metrics.PacketsAny.Load(),
		PacketsErr:        metrics.PacketsErr.Load(),
		LastSender:        lastSender,
		LastPacketTime:    lastPacketTime,
		LastNonJSONHead:   lastNonJSONHead,
		LastStatus:        lastStatus,
		LastPosition:      metrics.LastPosition.Load(),
		LastFuel:          &lastFuel,
		LastFlightTime:    &lastFlightTime,
		LastDistance:      &lastDistance,
		UpdateFlightErr:   metrics.UpdateFlightErr.Load(),
		UpdatePositionErr: metrics.UpdatePositionErr.Load(),
	}
}
