package server

import (
	"context"
	"time"

	"github.com/DecxBase/gateway/registry"
)

func LaunchHealthCheck(ctx context.Context, sp ServerPool) {
	t := time.NewTicker(time.Second * 20)
	registry.Logger.Info().Msg("Starting health check...")
	for {
		select {
		case <-t.C:
			go HealthCheck(ctx, sp)
		case <-ctx.Done():
			registry.Logger.Info().Msg("Closing Health Check")
			return
		}
	}
}
