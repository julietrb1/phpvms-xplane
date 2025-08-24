package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"

	"github.com/julietrb1/phpvms-xplane/internal/config"
	"github.com/julietrb1/phpvms-xplane/internal/service"
	"github.com/julietrb1/phpvms-xplane/internal/udp"
)

func Run(ctx context.Context, cancel context.CancelFunc, metrics *udp.Metrics, flightService *service.FlightService, cfg *config.Config, logger *slog.Logger) error {
	model := NewModel(ctx, cancel, metrics, flightService, cfg, logger)
	p := tea.NewProgram(&model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		logger.Error("Error running TUI", "error", err)
		return err
	}

	return nil
}
