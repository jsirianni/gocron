package cmd


var cfgFile      string
var frontendPort string
var summary      bool
var verbose      bool

var config Config


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
      Smtpserver   string `yaml:"smtpserver"`
      Smtpport     string `yaml:"smtpport"`
      Smtpaddress  string `yaml:"smtpaddress"`
      Smtppassword string `yaml:"smtppassword"`
      Interval     int    `yaml:"interval"`
      SlackHookUrl string `yaml:"slackhookurl"`
      SlackChannel string `yaml:"slackchannel"`
      PreferSlack  bool   `yaml:"preferslack"`
}
