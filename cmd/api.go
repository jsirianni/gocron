package cmd
import (
	"github.com/spf13/cobra"
)


// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the gocron REST API",
	Run: func(cmd *cobra.Command, args []string) {
		startBackend()
	},
}


func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(apiCmd)
    rootCmd.PersistentFlags().StringVar(&apiPort, "port", "3000", "API listening port (defaults to 3000)")

}


func startapi() {
	gocron.Api(apiPort)
}
