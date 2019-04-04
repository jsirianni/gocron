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
	backendCmd.Flags().StringVar(&backendPort, "port", "3000", "Listening port (defaults to 3000)")

}


func startBackend() {
	gocron.StartBackend(backendPort)
}
