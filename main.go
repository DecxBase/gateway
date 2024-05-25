package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DecxBase/gateway/registry"
	"github.com/DecxBase/gateway/server"
)

func main() {
	logger := registry.InitLogger()
	defer logger.Slog()

	config, err := registry.GetLBConfig()
	if err != nil {
		registry.Logger.Fatal().Msg(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPool, err := server.NewServerPool(registry.GetLBStrategy(config.Strategy))
	if err != nil {
		registry.Logger.Fatal().Msg(err.Error())
	}
	loadBalancer := server.NewLoadBalancer(serverPool)

	for _, u := range config.Backends {
		endpoint, err := url.Parse(u)
		if err != nil {
			logger.Fatal().Str("url", u).Msg(err.Error())
		}

		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := registry.NewBackend(endpoint, rp)
		rp.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			logger.Error().Str("error", e.Error()).Str("host", endpoint.Host).Msg("error handling the request")
			backendServer.SetAlive(false)

			if !server.AllowRetry(request) {
				logger.Error().Str("address", request.RemoteAddr).Str("path", request.URL.Path).Msg("Max retry attempts reached, terminating")
				http.Error(writer, "Service not available/blocked", http.StatusServiceUnavailable)
				return
			}

			logger.Error().Str("address", request.RemoteAddr).Str("path", request.URL.Path).Bool("retry", true).Msg("Attempting retry")

			loadBalancer.Serve(
				writer,
				request.WithContext(
					context.WithValue(request.Context(), server.RETRY_ATTEMPTED, true),
				),
			)
		}

		serverPool.AddBackend(backendServer)
	}

	serv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(loadBalancer.Serve),
	}

	go server.LaunchHealthCheck(ctx, serverPool)

	go func() {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := serv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Info().Int("port", config.Port).Msg("Load Balancer started")
	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Info().Str("error", err.Error()).Msg("ListenAndServe() error")
	}
}
