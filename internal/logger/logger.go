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
	// Log будет доступен всему коду как синглтон.
	// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
	Log *zap.Logger = zap.NewNop()

	// Sugar *zap.SugaredLogger.
	// Предоставляет удобный интерфейс для логирования с форматированием.
	Sugar *zap.SugaredLogger = Log.Sugar()

	// config хранит конфигурацию логера.
	config zap.Config
)

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("Initialize: %w", err)
	}

	config = zap.Config{
		Level:            lvl,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// создаём логер на основе конфигурации
	zl, err := config.Build()
	if err != nil {
		return fmt.Errorf("Initialize: %w", err)
	}
	// устанавливаем синглтон
	Log = zl
	Sugar = Log.Sugar()

	return nil
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h(w, r)
	})
}
