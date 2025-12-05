package serverless

import (
	"os"
	"time"

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

const checkDeploymentStatusInterval = 1 * time.Second

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
	logger  logr.Logger
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
	if logger.IsZero() {
		setupLogger()
	}

	if isDebug {
		logger.Info("debug logs enabled")
	}

	version.Print(logger)

	logger.Info("Running in serverless mode...")

	if keepAlive {
		keepProcessAlive()
	}
	return nil
}

func keepProcessAlive() {
	logger.V(1).Info("Running in keep-alive mode")

	t := time.NewTicker(checkDeploymentStatusInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if requiredOneAgentVersionIsDeployed() {
				logger.Info("OneAgent has been successfully deployed")
				t.Stop()
			} else {
				logger.V(1).Info("The required OneAgent version has not been deployed yet")
			}
		}
	}
}

func requiredOneAgentVersionIsDeployed() bool {
	// TODO: Check whether the target directory contains required OneAgent version
	return false
}

func setupLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.StacktraceKey = "stacktrace"

	logLevel := zapcore.InfoLevel
	if isDebug {
		// zap's debug level is -1, however this is not a valid value for the logr.Logger, so we have to overrule it.
		// use logger.V(1).Info to create debug logs.
		logLevel = zap.DebugLevel
	}

	zapLog := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(config), os.Stdout, logLevel))
	logger = zapr.NewLogger(zapLog)
}

func SetLogger(log logr.Logger) {
	logger = log
}
