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
var config       libgocron.Config


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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/gocron/config.yml", "config file (default is /etc/gocron/config.yml")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable standard out output with --verbose (default is disabled)" )
}


// initConfig reads ENV variables
func initConfig() {
	var err error
	config.Dbdatabase = os.Getenv("GC_DBDATABASE")
	config.Dbfqdn = os.Getenv("GC_DBFQDN")
	config.Dbpass = os.Getenv("GC_DBPASS")
	config.Dbport = os.Getenv("GC_DBPORT")
	config.Dbuser = os.Getenv("GC_DBUSER")
	config.SlackChannel = os.Getenv("GC_SLACK_CHANNEL")
	config.SlackHookURL = os.Getenv("GC_SLACK_HOOK_URL")

	config.Interval, err = strconv.Atoi(os.Getenv("GC_INTERVAL"))
	if err != nil {
		libgocron.LogError(err)
		os.Exit(1)
	}

	err = config.Validate()
	if err != nil {
		libgocron.LogError(err)
		os.Exit(1)
	}
}
