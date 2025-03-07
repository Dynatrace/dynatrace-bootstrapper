package main

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugFlag          = "debug"
	SuppressErrorsFlag = "suppress-error"
)

var (
	log                 logr.Logger
	isDebug             bool
	areErrorsSuppressed bool
)

func main() {
	cmd := bootstrapper(afero.NewOsFs())

	err := cmd.Execute()
	if err != nil {
		log.Error(err, "Error during bootstrapping")
		os.Exit(1)
	}
}

func bootstrapper(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "dynatrace-bootstrapper",
		RunE:               run(fs),
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	}

	AddFlags(cmd)
	move.AddFlags(cmd)
	configure.AddFlags(cmd)

	return cmd
}

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&isDebug, DebugFlag, false, "(Optional) Enables debug logs.")
	cmd.PersistentFlags().Lookup(DebugFlag).NoOptDefVal = "true"

	cmd.PersistentFlags().BoolVar(&areErrorsSuppressed, SuppressErrorsFlag, false, "(Optional) Always return exit code 0, even on error")
	cmd.PersistentFlags().Lookup(SuppressErrorsFlag).NoOptDefVal = "true"
}

func run(fs afero.Fs) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		setupLogger()

		if isDebug {
			log.Info("debug logs enabled")
		}
		version.Print(log)

		aferoFs := afero.Afero{
			Fs: fs,
		}

		err := move.Execute(log, aferoFs, move.SourceFolder, move.TargetFolder)
		if err != nil {
			if areErrorsSuppressed {
				log.Error(err, "error during moving, the error was suppressed")

				return nil
			}
			return err
		}

		err = configure.Execute(log, aferoFs, move.TargetFolder)
		if err != nil {
			if areErrorsSuppressed {
				log.Error(err, "error during configuration, the error was suppressed")

				return nil
			}
			return err
		}

		return nil
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
