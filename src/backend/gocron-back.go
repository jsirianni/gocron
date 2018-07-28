package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jsirianni/gocronlib"
)

const (
	version    string = "3.0.0"
	libVersion string = gocronlib.Version
)

var (
	verbose    bool             // Command line flag
	getVersion bool             // Command line flag
	config     gocronlib.Config = gocronlib.GetConfig(verbose)
)


func main() {
	fmt.Println("this")
	flag.BoolVar(&getVersion, "version", false, "Get the version and then exit")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	if getVersion == true {
		fmt.Println("gocron-back version:", version)
		fmt.Println("gocronlib version:", libVersion)
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
