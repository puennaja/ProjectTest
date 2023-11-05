package cmds

import (
	"github.com/spf13/cobra"
)

/*
	|--------------------------------------------------------------------------
	| Application's Command
	|--------------------------------------------------------------------------
	|
	| Here is which command you want to provide for your application
	| to use in your application.
	|
*/

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "user task",
	Short:   "task manages user task",
	Version: "v1.0.0",
}

func init() {
	rootCmd.AddCommand(serveRESTCmd)
}
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
