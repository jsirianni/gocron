package cmd
import (
	"gocron/libgocron"

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
	backendCmd.Flags().BoolVar(&summary, "summary", false, "Get summary")
}


func startBackend() {
	if summary == true {
		libgocron.GetSummary(verbose)
		return
	}

	libgocron.StartBackend(verbose)
}
