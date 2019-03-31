package libgocron
import (
    "testing"
)

func BuildTestConfig() Config {
    var c Config

    c.Dbdatabase = "test"
    c.Dbfqdn = "localhost"
    c.Dbpass = "password"
    c.Dbport = "5234"
    c.Dbuser = "test"
    c.Interval = 5
    c.PreferSlack = false
    c.SlackChannel = "test"
    c.SlackHookUrl = "http://badurl.com"

    return c
}

func BuildTestCron() Cron {
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

func Testalert(t *testing.T) {

    x := alert(BuildTestCron(), "test", "test")
    if x == false {
         t.Errorf("Expected alert to return true")
    }
}
