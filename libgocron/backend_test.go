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
    g.SlackHookURL = ""

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

func TestAlert(t *testing.T) {
    g := getTestConfig()
    c := getTestCron()


    // test bad hook url
    g.SlackHookURL = "http://badurl.com"
    if g.alert(c, "test", "test") == true {
        t.Errorf("Expected alert() to return false, due to bad Gocron config")
    }

    // should return true if 200 ok
    g.SlackHookURL = "https://httpstat.us/200"
    if g.alert(c, "test", "test") == false {
        t.Errorf("Expected alert() to return true, using mock http server")
    }
}

func TestslackAlert(t *testing.T) {
    g := getTestConfig()

    g.SlackHookURL = "http://badurl.com"
    if err := g.slackAlert("test", "test"); err == nil {
        t.Errorf("Expected slackAlert() to return an error when using a bad url")
    }

    g.SlackHookURL = "https://httpstat.us/200"
    if err := g.slackAlert("test", "test"); err != nil {
        t.Errorf("Expected slackAlert() to return a nil error, instead got: " + err.Error())
    }
}
