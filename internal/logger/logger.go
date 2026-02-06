// Package logger provides structured logging using go.uber.org/zap for the GophKeeper application.
// It offers a singleton Sugar logger instance for convenient logging throughout the application.
//
// Usage:
//
//	import "github.com/Pklerik/gophKeep/internal/logger"
//
//	func init() {
//		logger.InitLogger("info")
//	}
//
//	func main() {
//		logger.Sugar.Info("Application started")
//		logger.Sugar.Errorf("Error occurred: %v", err)
//	}
//
// The logger provides methods for different log levels:
// - Debug, Debugf, Debugw
// - Info, Infof, Infow
// - Warn, Warnf, Warnw
// - Error, Errorf, Errorw
// - Fatal, Fatalf, Fatalw
package logger

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var (
	// Log will be available across the codebase as a singleton.
	// By default a no-op logger is set which does not output any messages.
	Log *zap.Logger = zap.NewNop()

	// Sugar is a *zap.SugaredLogger providing a convenient formatted logging interface.
	Sugar *zap.SugaredLogger = Log.Sugar()
)

// Initialize initializes the singleton logger with the specified logging level.
func Initialize(config zap.Config) error {
	// build logger from configuration
	zl, err := config.Build()
	if err != nil {
		return fmt.Errorf("Initialize: %w", err)
	}
	// set the singleton
	Log = zl
	Sugar = Log.Sugar()

	return nil
}

// RequestLogger is a middleware logger for incoming HTTP requests.
func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h(w, r)
	})
}
