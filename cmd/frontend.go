package cmd
import (
	"github.com/spf13/cobra"
)


// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Start the frontend server",
	Long:  "Start the gocron frontend server, which presents an API that supports GET and POST requests",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}


func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(frontendCmd)
	frontendCmd.Flags().StringVar(&frontendPort, "port", "8080", "Listening port (defaults to 8080)")
}

func start() {
	gocron.StartFrontend(frontendPort)
}
