package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
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

VerboseRunLevel, DebugRunLevel bool
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&VerboseRunLevel, "verbose",
		"v", false, "turn on info logging")
	RootCmd.PersistentFlags().BoolVarP(&DebugRunLevel,"debug",
		"d", false, "turn on debug logging")
}
