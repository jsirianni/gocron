package libgocron

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
