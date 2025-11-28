package serverless

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Use = "serverless"

	TargetFolderFlag = "target"
	KeepAliveFlag    = "keep-alive"
	TechnologyFlag   = "technology"
	DebugFlag        = "debug"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:                Use,
		RunE:               run,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Version:            version.Version,
		Short:              "Deploy the OneAgent CodeModule in a Cloud environment",
	}

	addFlags(cmd)

	return cmd
}

var (
	log     logr.Logger
	isDebug bool

	targetFolder string
	technology   string
	keepAlive    bool
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&targetFolder, TargetFolderFlag, "", "Base path where to copy the CodeModule to.")
	err := cmd.MarkFlagRequired(TargetFolderFlag)
	if err != nil {
		panic(err)
	}

	cmd.Flags().BoolVar(&keepAlive, KeepAliveFlag, false, "Keep the Bootstrapper process running even after deployment is finished.")
	err = cmd.MarkFlagRequired(KeepAliveFlag)
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&technology, TechnologyFlag, "", "(Optional) Comma-separated list of CodeModule technologies to deploy.")
	cmd.Flags().BoolVar(&isDebug, DebugFlag, false, "(Optional) Enables debug logs.")
}

func run(_ *cobra.Command, _ []string) error {
	setupLogger()

	if isDebug {
		log.Info("debug logs enabled")
	}

	version.Print(log)

	log.Info("Running in serverless mode...")
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
