package gin

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/http/gin/ratelimiter"

	"github.com/chenjiandongx/ginprom"
	inspector "github.com/fatihkahveci/gin-inspector"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// MiddlewaresConfig is a set of middlewares related config params
type MiddlewaresConfig struct {
	DebugEnabled      bool
	PrometheusEnabled bool
	RateLimiterConfig struct {
		Enabled    bool
		Interval   time.Duration
		BucketSize int
	}
	SecureOptions struct {
		Enabled      bool
		AllowedHosts []string
		SSLHost      string
	}
	CorsOptions struct {
		Enabled         bool
		AllowOrigins    []string
		AllowMethods    []string
		AllowHeaders    []string
		ExposeHeader    []string
		AllowOriginFunc func(origin string) bool
		MaxAge          time.Duration
	}
	StaticFilesOptions struct {
		Enabled    bool
		ServeFiles []struct {
			Prefix                 string
			FilePath               string
			AllowDirectoryIndexing bool
		}
	}
	NewRelicOptions struct {
		ServiceName string
		LicenseKey  string
	}
}

// AddBasicHandlers will add basic handlers required by a Server.
// Library Ref: https://github.com/gin-gonic/contrib
func AddBasicHandlers(router *gin.Engine, config *MiddlewaresConfig, logger *logger.Logger) error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.NewRelicOptions.ServiceName),
		newrelic.ConfigLicense(config.NewRelicOptions.LicenseKey),
		// newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
	)
	if nil != err {
		return fmt.Errorf("new relic app initialisation: %w", err)
	}
	router.Use(nrgin.Middleware(app))

	// Setup the Health check middleware at the top
	router.Use(healthcheck.Default())

	// Documentation Swagger
	{
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Setup the Static file middleware handler before the inspector
	if config.StaticFilesOptions.Enabled {
		// if Allow DirectoryIndex
		// r.Use(static.Serve("/", static.LocalFile("/tmp", true)))
		// set prefix
		// r.Use(static.Serve("/static", static.LocalFile("/tmp", true)))
		for _, serveOption := range config.StaticFilesOptions.ServeFiles {
			router.Use(static.Serve(serveOption.Prefix, static.LocalFile(serveOption.FilePath, serveOption.AllowDirectoryIndexing)))
		}
	}

	// Setup secure options handler
	if config.SecureOptions.Enabled {
		router.Use(secure.Secure(secure.Options{
			AllowedHosts:          config.SecureOptions.AllowedHosts,
			SSLHost:               config.SecureOptions.SSLHost,
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:            315360000,
			SSLRedirect:           true,
			STSIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
		}))
	}

	// Setup the Cors support handler
	if config.CorsOptions.Enabled {
		// Setup CORS support
		// CORS for https://foo.com and https://github.com origins, allowing:
		// - PUT and PATCH methods
		// - Origin header
		// - Credentials share
		// - Preflight requests cached for 12 hours
		router.Use(cors.New(cors.Config{
			AllowOrigins:     config.CorsOptions.AllowOrigins,
			AllowMethods:     config.CorsOptions.AllowMethods,
			AllowHeaders:     config.CorsOptions.AllowHeaders,
			ExposeHeaders:    config.CorsOptions.ExposeHeader,
			AllowCredentials: true,
			AllowOriginFunc:  config.CorsOptions.AllowOriginFunc,
			MaxAge:           config.CorsOptions.MaxAge,
		}))
	}

	// Setup Debug inspector handler before the API handlers and after the static file handlers
	if config.DebugEnabled {
		router.Delims("{{", "}}")
		router.SetFuncMap(template.FuncMap{
			"formatDate": formatDate,
		})

		router.LoadHTMLFiles("assets/inspector/inspector.html")
		router.Use(inspector.InspectorStats())
		router.GET("/_inspector", func(c *gin.Context) {
			c.HTML(http.StatusOK, "inspector.html", map[string]interface{}{
				"title":      "Request Inspector",
				"pagination": inspector.GetPaginator(),
			})
		})
	}

	// Setup ratelimiter handler for user requests
	if config.RateLimiterConfig.Enabled {
		router.Use(ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
			// config: how to generate a limit key
			LimitKey: func(c *gin.Context) string {
				return c.ClientIP()
			},
			// config: how to respond when limiting
			LimitedHandler: func(c *gin.Context) {
				c.JSON(200, "Exceeds Request rate limit!!!")
				c.Abort()
			},
			// config: return ratelimiter token fill interval and bucket size
			TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
				return config.RateLimiterConfig.Interval, config.RateLimiterConfig.BucketSize
			},
		}, logger))
	}

	// Setup Prometheus exporter handler
	if config.PrometheusEnabled {
		// use prometheus metrics exporter middleware.
		//
		// ginprom.PromMiddleware() expects a ginprom.PromOpts{} poniter.
		// It is used for filtering labels by regex. `nil` will pass every requests.
		//
		// ginprom promethues-labels:
		//   `status`, `endpoint`, `method`
		//
		// for example:
		// 1). I don't want to record the 404 status request. That's easy for it.
		// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexStatus: "404"})
		//
		// 2). And I wish to ignore endpoints started with `/prefix`.
		// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexEndpoint: "^/prefix"})
		router.Use(ginprom.PromMiddleware(nil))

		// register the `/metrics` route.
		router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	}

	// Setup Debug PPROF handler before the API handlers
	if config.DebugEnabled {
		pprof.Register(router)
	}

	// Setup Error handler
	// router.Use(errors.Handler)
	return nil
}

func formatDate(t time.Time) string {
	return t.Format(time.RFC822)
}
