package k8sinit

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit/configure"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Use = "k8s-init"

	SourceFolderFlag   = "source"
	TargetFolderFlag   = "target"
	DebugFlag          = "debug"
	SuppressErrorsFlag = "suppress-error"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:                Use,
		RunE:               RunE,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Version:            version.Version,
		Short:              "Deploy the OneAgent CodeModule in a Kubernetes environment",
	}

	AddFlags(cmd)

	return cmd
}

var (
	log                 logr.Logger
	isDebug             bool
	areErrorsSuppressed bool

	sourceFolder string
	targetFolder string
)

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&sourceFolder, SourceFolderFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkFlagRequired(SourceFolderFlag)

	cmd.Flags().StringVar(&targetFolder, TargetFolderFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkFlagRequired(TargetFolderFlag)

	cmd.Flags().BoolVar(&isDebug, DebugFlag, false, "(Optional) Enables debug logs.")

	cmd.Flags().Lookup(DebugFlag).NoOptDefVal = "true"

	cmd.Flags().BoolVar(&areErrorsSuppressed, SuppressErrorsFlag, false, "(Optional) Always return exit code 0, even on error")

	cmd.Flags().Lookup(SuppressErrorsFlag).NoOptDefVal = "true"

	move.AddFlags(cmd)
	configure.AddFlags(cmd)
}

func RunE(_ *cobra.Command, _ []string) error {
	setupLogger()

	if isDebug {
		log.Info("debug logs enabled")
	}

	version.Print(log)

	err := move.Execute(log, sourceFolder, targetFolder)
	if err != nil {
		if areErrorsSuppressed {
			log.Error(err, "error during moving, the error was suppressed")

			return nil
		}

		log.Error(err, "error during configuration")

		return err
	}

	err = configure.SetupOneAgent(log, targetFolder)
	if err != nil {
		if areErrorsSuppressed {
			log.Error(err, "error during oneagent setup, the error was suppressed")

			return nil
		}

		log.Error(err, "error during configuration")

		return err
	}

	err = configure.EnrichWithMetadata(log)
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
