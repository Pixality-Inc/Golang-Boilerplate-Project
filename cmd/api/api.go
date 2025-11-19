package main

import "github.com/pixality-inc/golang-boilerplate-project/internal/wiring"

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	go func() {
		if err := wire.StartHealthcheckServer(); err != nil {
			wire.Log.WithError(err).Fatal("failed to start healthcheck server")
		}
	}()

	go func() {
		if err := wire.StartMetricsServer(); err != nil {
			wire.Log.WithError(err).Fatal("failed to start metrics server")
		}
	}()

	go func() {
		if err := wire.StartApiServer(); err != nil {
			wire.Log.WithError(err).Fatal("failed to start http server")
		}
	}()

	wire.ControlFlow.WaitForInterrupt()
}
