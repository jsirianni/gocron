package cmd
import (
	"os"
	"strconv"

	"gocron/libgocron"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

// global CLI variables
var cfgFile      string
var frontendPort string
var summary      bool
var verbose      bool
var gocron       libgocron.Gocron


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gocron",
	Short: "Monitor uptime with gocron",
	Long: "Monitor uptime with gocron",
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		libgocron.LogError(err)
		os.Exit(1)
	}
}


func init() {
	cobra.OnInitialize(initConfig)
}


// initConfig reads ENV variables
func initConfig() {
	var err error
	gocron.Dbdatabase = os.Getenv("GC_DBDATABASE")
	gocron.Dbfqdn = os.Getenv("GC_DBFQDN")
	gocron.Dbpass = os.Getenv("GC_DBPASS")
	gocron.Dbport = os.Getenv("GC_DBPORT")
	gocron.Dbuser = os.Getenv("GC_DBUSER")
	gocron.SlackChannel = os.Getenv("GC_SLACKCHANNEL")
	gocron.SlackHookURL = os.Getenv("GC_SLACKHOOKURL")

	gocron.Interval, err = strconv.Atoi(os.Getenv("GC_INTERVAL"))
	if err != nil {
		libgocron.LogError(err)
		os.Exit(1)
	}

	err = gocron.Validate()
	if err != nil {
		libgocron.LogError(err)
		os.Exit(1)
	}
}
