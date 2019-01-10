package libgocron

const VERSION string = "5.0"

var config Config
var verbose bool

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


type Config struct {
      Dbfqdn       string `yaml:"dbfqdn"`
      Dbport       string `yaml:"dbport"`
      Dbuser       string `yaml:"dbuser"`
      Dbpass       string `yaml:"dbpass"`
      Dbdatabase   string `yaml:"dbdatabase"`
      Interval     int    `yaml:"interval"`
      SlackHookUrl string `yaml:"slackhookurl"`
      SlackChannel string `yaml:"slackchannel"`
      PreferSlack  bool   `yaml:"preferslack"`
}
