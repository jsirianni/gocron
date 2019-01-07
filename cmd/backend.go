package cmd
import (
	"os"

	"gocron/libgocron"

	"github.com/spf13/cobra"
)


// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Start the gocron backend server",
	Long: "Start the gocron backend server, which alerts on missed jobs",
	Run: func(cmd *cobra.Command, args []string) {
		startBackend()
	},
}


func init() {
	rootCmd.AddCommand(backendCmd)
	backendCmd.Flags().BoolVar(&summary, "summary", false, "Get summary")
}


func startBackend() {
	// Initilize the config struct
	//config = GetConfig(verbose)

	if summary == true {
		// If verbose == true, summary will send to syslog AND the configured
		// alert system
		libgocron.GetSummary(verbose)
		return
	}

	// create the gocron table, if not exists
	if libgocron.CreateGocronTable(verbose) == false {
		os.Exit(1)
	}

	libgocron.Timer(verbose)
}
