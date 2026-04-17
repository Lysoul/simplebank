package ginserver

import (
	"assignments/simplebank/adapters/monitoring"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	// "github.com/uptrace/opentelemetry-go-extra/otelzap"
	// "go.opentelemetry.io/otel/trace"

	// ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// we could have a stuct with these things attached to it
// but whats the point...
//
//nolint:gochecknoglobals // we want to have a global config.
var config Config

type Config struct {
	Mode    string   `envconfig:"HTTP_MODE" default:"debug"`
	LogSkip []string `envconfig:"HTTP_LOG_SKIP" default:"/health,/metrics"`
	Port    int      `envconfig:"HTTP_PORT" default:"3000"`
	Prefix  string   `envconfig:"HTTP_PATH_PREFIX" default:""`

	EnableCORS          bool          `envconfig:"HTTP_ENABLE_CORS" default:"false"`
	AllowedOrigins      []string      `envconfig:"HTTP_ALLOWED_ORIGINS" default:"*"`
	AllowedHeaders      []string      `envconfig:"HTTP_ALLOWED_HEADERS" default:"*"`
	CORSMaxAge          time.Duration `envconfig:"HTTP_CORS_MAX_AGE" default:"24h"`
	DefaultCacheControl string        `envconfig:"HTTP_DEFAULT_CACHE_CONTROL" default:"no-store"`

	// generate etag for all responses
	EnableEtag bool `envconfig:"HTTP_ENABLE_ETAG" default:"false"`

	// default value to put in cache control, leave empty to disable
	// `private, no-store, no-cache` are good options
	CacheContolDefault string `envconfig:"HTTP_CACHE_CONTROL_DEFAULT"`

	ServiceName string `envconfig:"SERVICE_NAME" default:"service"`
}

// returns router and function to run server.
func InitGin(_config Config) (*gin.Engine, func() (*http.Server, func(context.Context))) {
	config = _config

	if config.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.ContextWithFallback = true

	if config.EnableCORS {
		monitoring.Logger().Info("CORS enabled")
		router.Use(CORSMiddleware(config))
	} else {
		monitoring.Logger().Info("CORS disabled")
	}

	return router, func() (*http.Server, func(context.Context)) {
		return Run(router)
	}
}

// to do return servver AND shutdown.
func Run(router *gin.Engine) (*http.Server, func(context.Context)) {
	var handler http.Handler
	if config.EnableEtag {
		handler = EtagHandler(router)
	} else {
		handler = router
	}
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.L().Info("ListenAndServe", zap.Error(err))
		}
	}()

	return server, func(ctx context.Context) {
		err := server.Shutdown(ctx)
		if err != nil {
			zap.L().Error("Server forced to shutdown:", zap.Error(err))
		} else {
			zap.L().Info("HTTP server has shut down")
		}
	}
}
