package main

import (
	"account-service/log"
	"account-service/tracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var httpClient *tracing.HTTPClient

func main() {
	zapOptions := []zap.Option{
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
		zap.IncreaseLevel(zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l != zapcore.DebugLevel })),
	}

	logger, _ := zap.NewDevelopment(zapOptions...)
	zapLogger := logger.With(zap.String("service", "account-service"))
	factory := log.NewFactory(zapLogger)
	// Init OTEL
	tp := tracing.TracerProvider("account-service", "otlp", factory)

	httpClient = tracing.NewHTTPClient(tp)

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware("account-service", otelecho.WithTracerProvider(tp)))

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type Message struct {
	Message string `json:"message"`
}

// Handler
func hello(c echo.Context) error {
	var dummy Message

	err := httpClient.GetJSON(c.Request().Context(), "", "http://localhost:1324/", &dummy)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dummy)
}
