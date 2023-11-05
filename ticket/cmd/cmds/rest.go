package cmds

import (
	"ticket/protocol"

	"github.com/spf13/cobra"
)

var serveRESTCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "start a http server",
	RunE:  serveREST,
}

// @title Ticket API
// @version 1.0

// @schemes http https
// @BasePath /
func serveREST(cmd *cobra.Command, args []string) error {
	return protocol.ServeREST()
}
