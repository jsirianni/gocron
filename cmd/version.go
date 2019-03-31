package cmd
import (
	"fmt"

	"gocron/libgocron"

	"github.com/spf13/cobra"
)


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the gocron version",
	Run: func(cmd *cobra.Command, args []string) {
		getVersion()
	},
}


func init() {
	rootCmd.AddCommand(versionCmd)
}


func getVersion() {
	fmt.Println("gocron version:", libgocron.Version)
}
