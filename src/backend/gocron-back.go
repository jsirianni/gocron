package main

import (
	"flag"
	"fmt"
	"time"

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
	config     gocronlib.Config = gocronlib.GetConfig(verbose)
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

	//if summary == true {
	//	gocronlib.CronLog("gocron retrieving summary", verbose)

		// If verbose == true, summary will send to syslog AND the configured
		// alert system
	//	getSummary()
	//	return
	//}

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
