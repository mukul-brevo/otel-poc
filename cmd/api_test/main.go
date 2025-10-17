package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/DTSL/golang-libraries/ctxhttpsrv"
	"github.com/DTSL/golang-libraries/di"
	"github.com/DTSL/golang-libraries/envutils"
	"github.com/DTSL/golang-libraries/errors"
	"github.com/DTSL/golang-libraries/mainutils"
	"go.uber.org/dig"
)

const (
	appname = "otel-poc-api"
	version = "1.0.0"
)

func main() {
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http/protobuf")
	// export OTEL_EXPORTER_OTLP_METRICS_PROTOCOL="http/protobuf"

	mainutils.Run(run)
}

func run(ctx context.Context) error {
	closeInit, err := mainutils.Init(mainutils.Config{
		AppName: appname,
		Version: version,
		Env:     envutils.Testing,
	})
	if err != nil {
		return errors.Wrap(err, "init")
	}
	defer closeInit()

	// Initialize DI container
	container, err := getContainer(ctx, di.Application{
		Env:     envutils.Testing,
		Name:    appname,
		Version: version,
	})
	if err != nil {
		return errors.Wrap(err, "get container")
	}

	// Close resources
	defer func() {
		errCloser := container.Invoke(func(closer *di.Closer) {
			_ = closer.Close()
		})
		if errCloser != nil {
			err = errors.Join(err, errCloser)
		}
	}()

	// Start server
	err = runHTTPServer(container, ":8088")
	if err != nil {
		return errors.Wrap(err, "start server")
	}
	return nil
}

func runHTTPServer(container *dig.Container, addr string) error {
	// Start server
	log.Printf("Start HTTP server on %s", addr)
	err := container.Invoke(func(ctx context.Context, handler http.Handler) error {
		return ctxhttpsrv.ListenAndServe(ctx, addr, handler)
	})
	if err != nil {
		return errors.Wrap(err, "ListenAndServe")
	}
	log.Println("Stopped HTTP server")
	return nil
}
