package cmd

import (
	"github.com/id27182/ami-proxy/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Start() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start ami proxy",
		Run: func(cmd *cobra.Command, args []string) {
			err := server.Serve()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
