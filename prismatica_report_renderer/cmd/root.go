package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"

	"github.com/Project-Prismatica/prismatica_report_renderer/configuration"
)

const (
	CommandFlagNameVerbose = "verbose"
	CommandFlagNameDebug = "debug"
	CommandFlagNameMongodbUri = "mongodb-uri"
)

var (

	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "prismatica_report_renderer",
		Short: "A gRPC based service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if DebugRunLevel {
				logrus.SetLevel(logrus.DebugLevel)
			} else if VerboseRunLevel {
				logrus.SetLevel(logrus.InfoLevel)
			}
		},
	}

	VerboseRunLevel = false
	DebugRunLevel = false
	MongodbUri string
)

func Execute(configurationSource *viper.Viper) {

	configurationSource.BindPFlag(configuration.ConfigKeyLogLevelDebug,
		RootCmd.Flags().Lookup(CommandFlagNameDebug))
	configurationSource.BindPFlag(configuration.ConfigKeyLogLevelVerbose,
		RootCmd.Flags().Lookup(CommandFlagNameVerbose))
	configurationSource.BindPFlag(configuration.ConfigKeyMongodbUri,
		RootCmd.Flags().Lookup(CommandFlagNameMongodbUri))

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&VerboseRunLevel, CommandFlagNameVerbose,
	"v", false, "turn on info logging")
	RootCmd.PersistentFlags().BoolVarP(&DebugRunLevel, CommandFlagNameDebug,
	"d", false, "turn on debug logging")
	RootCmd.PersistentFlags().StringVarP(&MongodbUri, CommandFlagNameMongodbUri,
	"u", "",
	"URI for default MongoDB template collection in the form of " +
			"mongodb://[[user:password]@]host:port/database/collection")
}
