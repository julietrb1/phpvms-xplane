package udp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type Listener struct {
	Addr     *net.UDPAddr
	Conn     *net.UDPConn
	Metrics  *Metrics
	Logger   *slog.Logger
	Handler  PayloadHandler
	MaxBytes int
}

type PayloadHandler interface {
	HandlePayload(ctx context.Context, payload *Payload) (error, error)
}

func NewListener(bindHost string, bindPort int, handler PayloadHandler, logger *slog.Logger) (*Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &Listener{
		Addr:     addr,
		Metrics:  NewMetrics(),
		Logger:   logger,
		Handler:  handler,
		MaxBytes: 64 * 1024, // 64 KiB max datagram size
	}, nil
}

func (l *Listener) Start(ctx context.Context) error {
	var err error
	l.Conn, err = net.ListenUDP("udp", l.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}
	defer l.Conn.Close()

	l.Logger.Info("UDP listener started", "addr", l.Addr.String())

	buffer := make([]byte, l.MaxBytes)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := l.Conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
				l.Logger.Warn("Failed to set read deadline", "error", err)
			}

			n, addr, err := l.Conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				l.Logger.Error("Error reading from UDP", "error", err)
				continue
			}

			l.processPacket(ctx, buffer[:n], addr)
		}
	}
}

func (l *Listener) processPacket(ctx context.Context, data []byte, addr *net.UDPAddr) {
	l.Metrics.PacketsAny.Add(1)
	l.Metrics.LastSender.Store(addr)
	l.Metrics.LastPacketTime.Store(time.Now().Unix())

	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		l.Metrics.PacketsErr.Add(1)
		// Store first 8 bytes of non-JSON data for debugging
		if len(data) > 0 {
			head := fmt.Sprintf("%x", data[:min(8, len(data))])
			l.Metrics.LastNonJSONHead.Store(&head)
		}
		l.Logger.Debug("Failed to decode JSON", "error", err, "addr", addr.String())
		return
	}

	l.Metrics.LastStatus.Store(&payload.Status)
	l.Metrics.LastPosition.Store(&payload.Position)
	l.Metrics.LastDistance.Store(int32(payload.Position.DistanceNM))
	l.Metrics.LastFuel.Store(int32(payload.Fuel))
	l.Metrics.LastFlightTime.Store(int32(payload.FlightTime))

	if l.Handler == nil {
		err := fmt.Errorf("no handler set")
		l.Metrics.UpdateFlightErr.Store(&err)
		l.Metrics.UpdatePositionErr.Store(&err)
		return
	}

	updateFlightErr, updatePositionErr := l.Handler.HandlePayload(ctx, &payload)
	l.Metrics.UpdateFlightErr.Store(&updateFlightErr)
	l.Metrics.UpdatePositionErr.Store(&updatePositionErr)
}

func (l *Listener) GetMetrics() *Metrics {
	return l.Metrics
}
