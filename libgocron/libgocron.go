package libgocron

import (
    "errors"
    "net/url"
)

// Version holds the current version of gocron
const Version string = "5.1.0"

// Cron defines a cronjob
type Cron struct {
      Cronname    string `json:cronname`  // Name of the cronjob
      Account     string `json:account`   // Account the job belongs to
      Email       string `json:email`     // Address to send alerts to
      Frequency   int    `json:frequency` // How often a job should check in
      Site        bool   `json:site`      // Set true if service is a site (Example: Network gateway)
      Ipaddress   string   // Source IP address
      Lastruntime int      // Unix timestamp
      Alerted     bool     // set to true if an alert has already been thrown
}

// Gocron defines a global configuration used at runtime
type Gocron struct {
      Dbfqdn       string `yaml:"dbfqdn"`
      Dbport       string `yaml:"dbport"`
      Dbuser       string `yaml:"dbuser"`
      Dbpass       string `yaml:"dbpass"`
      Dbdatabase   string `yaml:"dbdatabase"`
      Interval     int    `yaml:"interval"`
      SlackHookURL string `yaml:"slackhookurl"`
      SlackChannel string `yaml:"slackchannel"`
}

// Validate checks if config parameters are valid
func (g Gocron) Validate() error {
    message := "Errors found in the configuration:\n"
    m := ""

    if len(g.Dbdatabase) == 0 {
        m = m + "dbdatabase is length 0\n"
    }

    if len(g.Dbfqdn) == 0 {
        m = m + "dbfqdn is length 0\n"
    }

    if len(g.Dbpass) == 0 {
        m = m + "dbpass is length 0\n"
    }

    if len(g.Dbport) == 0 {
        m = m + "dbport is length 0\n"
    }

    if len(g.Dbuser) == 0 {
        m = m + "dbuser is length 0\n"
    }

    if g.Interval < 1 {
        m = m + "interval is less than 1\n"
    }

    if len(g.SlackChannel) == 0 {
        m = m + "slack_channel is length 0\n"
    }

    if len(g.SlackHookURL) == 0 {
        m = m + "slack_hook_url is length 0\n"
    } else {
        _, err := url.Parse(g.SlackHookURL)
        if err != nil {
            m = m + "slack_hook_url: "
            m = m + err.Error()
        }
    }

    if len(m) > 0 {
        return errors.New(message + m)
    }
    return nil
}
