package azureappservice

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Use = "azure-app-service"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:                Use,
		RunE:               run,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Version:            version.Version,
		Short:              fmt.Sprintf("%s version %s", version.AppName, version.Version),
	}

	return cmd
}

var (
	log     logr.Logger
	isDebug bool
)

func run(_ *cobra.Command, _ []string) error {
	setupLogger()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Info("Got interrupted, shutting down")

			// do cleanup...

			return nil
		default:
			sync.OnceFunc(func() {
				log.Info("I am serverless")

				// wait for lock file

				// do work
			})()

			time.Sleep(time.Second)
			log.Info("I am alive !!!!!!")
		}
	}
}

func setupLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.StacktraceKey = "stacktrace"

	logLevel := zapcore.InfoLevel
	if isDebug {
		// zap's debug level is -1, however this is not a valid value for the logr.Logger, so we have to overrule it.
		// use log.V(1).Info to create debug logs.
		logLevel = zap.DebugLevel
	}

	zapLog := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(config), os.Stdout, logLevel))

	log = zapr.NewLogger(zapLog)
}
