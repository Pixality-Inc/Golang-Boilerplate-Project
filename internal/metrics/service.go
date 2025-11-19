package metrics

import (
	"context"
	"time"

	"github.com/pixality-inc/golang-core/env"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/metrics"
)

type AppMetrics interface {
	SetUptime(duration time.Duration)
}

type BooksMetrics interface {
	IncBooksCalls()
}

type Service interface {
	Register(ctx context.Context) error
	Start(ctx context.Context)
}

type Impl struct {
	manager         metrics.Manager
	appEnv          env.AppEnv
	refreshDuration time.Duration
	uptime          metrics.Histogram
	booksCalls      metrics.Counter
}

func New(
	manager metrics.Manager,
	appEnv env.AppEnv,
	refreshDuration time.Duration,
) *Impl {
	return &Impl{
		manager:         manager,
		appEnv:          appEnv,
		refreshDuration: refreshDuration,
		booksCalls: manager.NewCounter(
			metrics.NewMetricDescription("books.calls"),
		),
		uptime: manager.NewHistogram(
			metrics.NewMetricDescription("app.uptime"),
			metrics.NewHistogramOptions(),
		),
	}
}

func (s *Impl) Register(ctx context.Context) error {
	logger.GetLogger(ctx).Info("Registering metrics...")

	return s.manager.Register(
		s.booksCalls,
		s.uptime,
	)
}

func (s *Impl) Start(ctx context.Context) {
	logger.GetLogger(ctx).Info("Starting metrics routine...")

	ticker := time.NewTicker(s.refreshDuration)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger(ctx).Info("Stopping metrics routine...")

			return

		case <-ticker.C:
			s.refreshMetrics()
		}
	}
}

func (s *Impl) IncBooksCalls() {
	s.booksCalls.Inc()
}

func (s *Impl) SetUptime(duration time.Duration) {
	s.uptime.Observe(float64(duration.Milliseconds()))
}

func (s *Impl) refreshMetrics() {
	s.SetUptime(time.Since(s.appEnv.StartedAt()))
}
