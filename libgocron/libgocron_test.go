package libgocron
import (
    "testing"
)

func getTestConfig() Gocron {
    var g Gocron

    g.Dbdatabase = "test"
    g.Dbfqdn = "localhost"
    g.Dbpass = "password"
    g.Dbport = "5234"
    g.Dbuser = "test"
    g.Interval = 5
    g.SlackChannel = "test"
    g.SlackHookURL = "http://valid.com"

    return g
}

func getTestCron() Cron {
    var c Cron

    c.Account = "test"
    c.Alerted = false
    c.Cronname = "test"
    c.Email = "test@test.com"
    c.Frequency = 20
    c.Ipaddress = "10.0.0.1"
    c.Lastruntime = 000
    c.Site = false

    return c
}

func TestValidate(t *testing.T) {
    g := getTestConfig()
    if err := g.Validate(); err != nil {
        t.Errorf("Expected Validate() to return nil when using a valid config")
    }

    // Invalid parameters
    g = getTestConfig()
    g.Dbdatabase = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Dbdatabase")
    }

    g = getTestConfig()
    g.Dbfqdn = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Dbfqdn")
    }

    g = getTestConfig()
    g.Dbpass = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Dbpass")
    }

    g = getTestConfig()
    g.Dbport = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Dbport")
    }

    g = getTestConfig()
    g.Dbuser = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Dbuser")
    }

    g = getTestConfig()
    g.Interval = -1
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad Interval")
    }

    g = getTestConfig()
    g.SlackChannel = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad SlackChannel")
    }

    g = getTestConfig()
    g.SlackChannel = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad SlackChannel")
    }

    g = getTestConfig()
    g.SlackHookURL = ""
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad SlackHookURL")
    }

    g.SlackHookURL = "4"
    if g.Validate() == nil {
        t.Errorf("Expected Validate() to return an error, when using bad SlackHookURL")
    }
}
