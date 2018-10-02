package main

import (
	"flag"
	"fmt"
	"time"
	"os"


	"../gocronlib"
)

const (
	version    string = "3.0.4"
	libVersion string = gocronlib.Version
)

var (
	summary    bool           // Command line flag
	verbose    bool           // Command line flag
	getVersion bool           // Command line flag
	config     gocronlib.Config
)


func main() {
	flag.BoolVar(&getVersion, "version", false, "Get the version and then exit")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&summary, "summary", false, "Enable weekly summary")
	flag.Parse()

	if getVersion == true {
		fmt.Println("gocron-back version:", version)
		fmt.Println("gocronlib version:", libVersion)
		return
	}

	// Build config
	config = gocronlib.GetConfig(verbose)

	if summary == true {
		// If verbose == true, summary will send to syslog AND the configured
		// alert system
		getSummary()
		return
	}

	if verbose == true {
		fmt.Println("Verbose mode enabled")
		fmt.Println("gocron-back version:", version)
		fmt.Println("gocronlib version:", libVersion)
		fmt.Println("Using check interval:", config.Interval)

		if config.PreferSlack == true {
			fmt.Println("Prefer slack: enabled")
			fmt.Println("Slack channel: " + config.SlackChannel)
			fmt.Println("Slack hook url: " + config.SlackHookUrl)

		}
	}

	// create the gocron table, if not exists
	if gocronlib.CreateGocronTable(verbose) == false {
		os.Exit(1)
	}


	timer()
}



// Function calls checkCronStatus() on a set interval
func timer() {
	for {
		time.Sleep((time.Duration(config.Interval) * time.Second))
		gocronlib.CronLog("Checking for missed jobs.", verbose)
		cronStatus()
	}
}
