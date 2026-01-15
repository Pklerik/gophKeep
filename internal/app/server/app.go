// Package server provides the server application startup and management.
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Pklerik/gophKeep/internal/config/server"
	"github.com/Pklerik/gophKeep/internal/logger"
	"github.com/Pklerik/gophKeep/internal/router/server"
	"golang.org/x/sync/errgroup"
)

// StartApp starts the server application.
func StartApp(flags config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		<-c
		cancel()
	}()

	routerHandler, err := server.NewServerRouter().ConfigureRouter(ctx, flags)
	if err != nil {
		logger.Sugar.Errorf("Unable to start server: %v", err)
		return
	}

	httpServer := &http.Server{
		Addr:         flags.ServerAddress,
		Handler:      routerHandler,
		ReadTimeout:  flags.Timeout,
		WriteTimeout: flags.Timeout,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if flags.TLS {
			return runTLSListener(httpServer, flags.CertFile, flags.KeyFile)
		}

		logger.Sugar.Infof("Starting server on %s", flags.ServerAddress)
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		logger.Sugar.Infof("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return httpServer.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		logger.Sugar.Errorf("Server error: %v", err)
	}
}

// runTLSListener runs the server with TLS/HTTPS.
func runTLSListener(httpServer *http.Server, certFile, keyFile string) error {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	httpServer.TLSConfig = tlsConfig

	logger.Sugar.Infof("Starting server with TLS on %s", httpServer.Addr)
	return httpServer.ListenAndServeTLS(certFile, keyFile)
}
