// Package cmd provides the command-li		PersistentPreRunE:  func(_ *cobra.Command, _ []string) error { return nil },e interface for the Dynatrace bootstrapper.
package cmd

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Use is the name of the CLI application.
	Use = "dynatrace-bootstrapper"

	// SourceFolderFlag is the flag for specifying the source folder.
	SourceFolderFlag = "source"
	// TargetFolderFlag is the flag for specifying the target folder.
	TargetFolderFlag = "target"
	// DebugFlag is the flag for enabling debug mode.
	DebugFlag = "debug"
	// SuppressErrorsFlag is the flag for suppressing errors.
	SuppressErrorsFlag = "suppress-error"
)

// New creates a new cobra command for the Dynatrace bootstrapper.
func New(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:               Use,
		RunE:              run(fs),
		Version:           version.Version,
		SilenceUsage:      true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return nil },
	}

	cmd.PersistentFlags().String(SourceFolderFlag, "", "source folder of the DT OneAgent installer")
	cmd.PersistentFlags().String(TargetFolderFlag, "", "target folder where the DT OneAgent should be copied to")
	cmd.PersistentFlags().Bool(DebugFlag, false, "enable debug output")
	cmd.PersistentFlags().Bool(SuppressErrorsFlag, false, "suppress one agent install errors")

	configure.AddFlags(cmd)
	move.AddFlags(cmd)

	return cmd
}

var (
	log                 logr.Logger
	isDebug             bool
	areErrorsSuppressed bool

	sourceFolder string
	targetFolder string
)

func run(fs afero.Fs) func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {
		setupLogger()

		if isDebug {
			log.Info("debug logs enabled")
		}

		version.Print(log)

		aferoFs := afero.Afero{
			Fs: fs,
		}

		err := move.Execute(log, aferoFs, sourceFolder, targetFolder)
		if err != nil {
			if areErrorsSuppressed {
				log.Error(err, "error during moving, the error was suppressed")

				return nil
			}

			log.Error(err, "error during configuration")

			return err
		}

		err = configure.SetupOneAgent(log, aferoFs, targetFolder)
		if err != nil {
			if areErrorsSuppressed {
				log.Error(err, "error during oneagent setup, the error was suppressed")

				return nil
			}

			log.Error(err, "error during configuration")

			return err
		}

		err = configure.EnrichWithMetadata(log, aferoFs)
		if err != nil {
			if areErrorsSuppressed {
				log.Error(err, "error during enrichment, the error was suppressed")

				return nil
			}

			log.Error(err, "error during enrichment")

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

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	isDebug, _ = rootCommand.PersistentFlags().GetBool(DebugFlag)
	areErrorsSuppressed, _ = rootCommand.PersistentFlags().GetBool(SuppressErrorsFlag)
	sourceFolder, _ = rootCommand.PersistentFlags().GetString(SourceFolderFlag)
	targetFolder, _ = rootCommand.PersistentFlags().GetString(TargetFolderFlag)
}

var rootCommand = New(afero.NewOsFs())
