package server

import (
	nethttp "net/http"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/config"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/http"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/tracer"
)

// DI is used to inject the server dependencies
type DI struct {
	Config                                        config.AppConfig
	Address                                       string
	Logger                                        *logger.Logger
	TracerProvider                                *tracer.Provider
	TracerProviderShutdownHandler, PreRunCallback func() error
	Handler                                       nethttp.Handler
}

// NewServer initialises the gRPC server.
func NewServer(di DI) (*http.Server, error) {
	return http.NewServer(
		&http.ServerConfig{
			Address: di.Address,
			GracefulShutdownHandler: func() error {
				return nil
			},
			Name:                          di.Config.ServiceName,
			TracerProvider:                di.TracerProvider,
			TracerProviderShutdownHandler: di.TracerProviderShutdownHandler,
		},
		di.Logger,
		di.PreRunCallback,
		di.Handler,
	)
}
