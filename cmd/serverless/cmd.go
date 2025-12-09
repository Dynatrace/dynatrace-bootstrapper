package serverless

import (
	"os"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/deployment"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/log"
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
	SourceFolderFlag = "source"
	DebugFlag        = "debug"
)

const (
	defaultCodeModulesPathInSourceFolder = "/opt/dynatrace/oneagent"
)

const checkDeploymentStatusInterval = 10 * time.Second

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

	sourceFolder string
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
	cmd.Flags().StringVar(&sourceFolder, SourceFolderFlag, defaultCodeModulesPathInSourceFolder, "(Optional) Base path where to copy the CodeModule from.")
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

	status, err := deployment.CheckAgentDeploymentStatus(sourceFolder, targetFolder)
	if err != nil {
		logger.Error(err, "failed to check OneAgent deployment status. Skipping deployment.", "status", status.String())
	} else if status == deployment.Deployed {
		logger.Info("OneAgent is already deployed")
	} else {
		logger.Info("OneAgent deployment status", "status", status)

		if status == deployment.NotDeployed {
			// TODO: Acquire the lock file
			// TODO: Deploy OneAgent code modules based on the new folder hierarchy
			// TODO: Create the active symlink
		} else if status == deployment.LinkMissing {
			// TODO: Acquire the lock
			// TODO: Create the active symlink
		}
	}

	if keepAlive {
		keepProcessAlive(status)
	}

	return err
}

func keepProcessAlive(status deployment.Status) {
	log.Debug(logger, "Running in keep-alive mode")

	deploymentStatusTicker := time.NewTicker(checkDeploymentStatusInterval)
	if status == deployment.Deployed {
		deploymentStatusTicker.Stop()
	}

	var lastErr error
	for {
		select {
		case <-deploymentStatusTicker.C:
			status, err := deployment.CheckAgentDeploymentStatus(sourceFolder, targetFolder)
			if err != nil {
				if lastErr == nil || err.Error() != lastErr.Error() {
					logger.Error(err, "failed to check OneAgent deployment status", "status", status.String())
					lastErr = err
				}

				continue
			}

			if status == deployment.Deployed {
				logger.Info("OneAgent has been successfully deployed")
				deploymentStatusTicker.Stop()
			} else {
				log.Debug(logger, "The required OneAgent version is not deployed", "status", status.String())
			}
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
		// use logger.V(1).Info to create debug logs.
		logLevel = zap.DebugLevel
	}

	zapLog := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(config), os.Stdout, logLevel))
	logger = zapr.NewLogger(zapLog)
}

func SetLogger(log logr.Logger) {
	logger = log
}
