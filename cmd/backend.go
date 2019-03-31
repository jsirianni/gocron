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
	gocron.StartBackend()
}
