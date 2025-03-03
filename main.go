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
	sourceFolderFlag = "source"
	targetFolderFlag = "target"
	debugFlag        = "debug"
	silentFlag       = "silent"
)

var (
	log      logr.Logger
	isDebug  bool
	isSilent bool

	sourceFolder string
	targetFolder string
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
		Use:  "dynatrace-bootstrapper",
		RunE: run(fs),
	}

	AddFlags(cmd)
	move.AddFlags(cmd)
	configure.AddFlags(cmd)

	return cmd
}

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&sourceFolder, sourceFolderFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkPersistentFlagRequired(sourceFolderFlag)

	cmd.PersistentFlags().StringVar(&targetFolder, targetFolderFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkPersistentFlagRequired(targetFolderFlag)

	cmd.PersistentFlags().BoolVar(&isDebug, debugFlag, false, "Enables debug logs.")
	cmd.PersistentFlags().BoolVar(&isSilent, silentFlag, false, "Always return exit code 0, even on error")
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

		err := move.Execute(log, aferoFs, sourceFolder, targetFolder)
		if err != nil {
			if isSilent {
				log.Error(err, "error during moving, silent fail is enabled")
				return nil
			}
			return err
		}

		err = configure.Execute(log, aferoFs, targetFolder)
		if err != nil {
			if isSilent {
				log.Error(err, "error during configuration, silent fail is enabled")
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
