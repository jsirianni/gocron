package cmd
import (
	"github.com/spf13/cobra"
)


// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Start the gocron backend server",
	Run: func(cmd *cobra.Command, args []string) {
		startBackend()
	},
}


func init() {
	rootCmd.AddCommand(backendCmd)
}


func startBackend() {
	// start the api on a new thread
	go gocron.Api(apiPort)

	// start the backend service
	gocron.StartBackend()
}
