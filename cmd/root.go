package cmd
import (
	"fmt"
	"os"
	"strconv"

	"gocron/libgocron"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println(err)
		os.Exit(1)
	}
}


func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/gocron/config.yml", "config file (default is /etc/gocron/config.yml")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable standard out output with --verbose (default is disabled)" )
}


// initConfig reads in config file and ENV variables
func initConfig() {

	// set the config file to be read
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// read the config file
	if err := viper.ReadInConfig(); err == nil {
		libgocron.CronLog("Starting gocron . . .")
		libgocron.CronLog("Using config file: " + viper.ConfigFileUsed())
	} else {
		libgocron.CronLog("Config file not found: " + cfgFile)
	}

	// read the environment variables
	viper.SetEnvPrefix("GC")
	viper.AutomaticEnv()

	// Unmarshal the configuration into config (Config struct)
	// environment values will replace values found in the config file
	//
	err := viper.Unmarshal(&config)
	if err != nil {
		libgocron.CronLog(err.Error())
		os.Exit(1)
	} else {
		libgocron.CronLog("Starting gocron with config: ")
		libgocron.CronLog("dbfqdn: " + config.Dbfqdn)
		libgocron.CronLog("dbport: " +  config.Dbport)
		libgocron.CronLog("dbuser: " +  config.Dbuser)
		libgocron.CronLog("dbdatabase: " +  config.Dbdatabase)
		libgocron.CronLog("interval: " +  strconv.Itoa(config.Interval))
		libgocron.CronLog("preferslack: " +  strconv.FormatBool(config.PreferSlack))
		libgocron.CronLog("slackchannel: " +  config.SlackChannel)
		libgocron.CronLog("slackhookurl: " +  config.SlackHookURL)
	}

	err = config.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}



	// TODO: implement this, which will likely require some logic that
	// /includes a dedicated "read config" function
	/*viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	read the config file into Config struct
	var c Config = GetConfig(verbose)
	*/
}
