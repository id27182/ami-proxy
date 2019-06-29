package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ami-proxy",
		Short: "Azure managed identity proxy",
		Long:  "",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debugFlag {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	debugFlag bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "v", false, "verbose output")
}

func Execute() {

	rootCmd.AddCommand(Start())

	rootCmd.Execute()
}
