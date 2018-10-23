package cmd
import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/gocron/config.yml", "config file (default is /etc/gocron/config.yml")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable standard out output with --verbose (default is disabled)" )

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables
func initConfig() {

	// set the config file to be read
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// read the config file
	if err := viper.ReadInConfig(); err == nil {
		CronLog("Starting gocron . . .", verbose)
		CronLog("Using config file: " + viper.ConfigFileUsed(), verbose)
	} else {
		CronLog("Config file not found: " + cfgFile, verbose)
	}

	// read the environment variables
	viper.SetEnvPrefix("GC")
	viper.AutomaticEnv()

	// Unmarshal the configuration into config (Config struct)
	// environment values will replace values found in the config file
	//
	err := viper.Unmarshal(&config)
	if err != nil {
		CronLog(err.Error(), verbose)
		os.Exit(1)
	} else {
		CronLog("Starting gocron with config: ", verbose)
		CronLog("dbfqdn: " + config.Dbfqdn, verbose)
		CronLog("dbport: " +  config.Dbport, verbose)
		CronLog("dbuser: " +  config.Dbuser, verbose)
		CronLog("dbdatabase: " +  config.Dbdatabase, verbose)
		CronLog("interval: " +  strconv.Itoa(config.Interval), verbose)
		CronLog("preferslack: " +  strconv.FormatBool(config.PreferSlack), verbose)
		CronLog("slackchannel: " +  config.SlackChannel, verbose)
		CronLog("slackhookurl: " +  config.SlackHookUrl, verbose)
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
