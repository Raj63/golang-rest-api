package config

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
)

// AppConfig maps the environment variables into a struct.
type AppConfig struct {
	// AppEnv is the application environment that determines `configs/<APP_ENV>.env` to load.
	AppEnv string `env:"APP_ENV" envDefault:"development"`

	// PrimaryDBUri is the primary database URI.
	PrimaryDBUri string `env:"PRIMARY_DB_URI"`

	// ServiceName is the application's service name.
	ServiceName string `env:"SERVICE_NAME" envDefault:"golang-rest-api"`

	// TracerCollectorAddress is the OpenTelemetry trace collector address.
	TracerCollectorAddress string `env:"TRACER_COLLECTOR_ADDRESS"`

	// NewRelicLicenseKey is the license key for New relic instrumentations
	NewRelicLicenseKey string `env:"NEWRELIC_LICENSE_KEY"`

	HTTPConfig struct {
		// Address is the HTTP server's address.
		Address string `env:"HTTP_ADDRESS"`

		// Enabled is the feature flag
		Enabled bool `env:"HTTP_ENABLED"`
	}

	HTTPSConfig struct {
		// Address is the HTTPS server's address.
		Address string `env:"HTTPS_ADDRESS"`

		// Enabled is the feature flag
		Enabled bool `env:"HTTPS_ENABLED"`
	}
}

// NewConfig loads <APP_ENV> into Config struct.
func NewConfig() (AppConfig, error) {
	config := AppConfig{}
	err := parseEnv(&config)
	return config, err
}

func parseEnv(c interface{}) error {
	if err := env.ParseWithFuncs(c, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]byte{}):            parseByteArray,
		reflect.TypeOf([][]byte{}):          parseByte2DArray,
		reflect.TypeOf(map[string]int{}):    parseMapStrInt,
		reflect.TypeOf(map[string]string{}): parseMapStrStr,
		reflect.TypeOf(http.SameSite(1)):    parseHTTPSameSite,
	}); err != nil {
		return err
	}

	return nil
}

func parseByteArray(v string) (interface{}, error) {
	return []byte(v), nil
}

func parseByte2DArray(v string) (interface{}, error) {
	newBytes := [][]byte{}
	bytes := strings.Split(v, ",")
	for _, b := range bytes {
		newBytes = append(newBytes, []byte(b))
	}

	return newBytes, nil
}

func parseHTTPSameSite(v string) (interface{}, error) {
	ss, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}

	return http.SameSite(ss), nil
}

func parseMapStrInt(v string) (interface{}, error) {
	newMaps := map[string]int{}
	maps := strings.Split(v, ",")
	for _, m := range maps {
		splits := strings.Split(m, ":")
		if len(splits) != 2 {
			continue
		}

		val, err := strconv.Atoi(splits[1])
		if err != nil {
			return nil, err
		}

		newMaps[splits[0]] = val
	}

	return newMaps, nil
}

func parseMapStrStr(v string) (interface{}, error) {
	newMaps := map[string]string{}
	maps := strings.Split(v, ",")
	for _, m := range maps {
		splits := strings.Split(m, ":")
		if len(splits) != 2 {
			continue
		}

		newMaps[splits[0]] = splits[1]
	}

	return newMaps, nil
}
