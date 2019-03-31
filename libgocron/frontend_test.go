package libgocron
import (
    "testing"
)

func TestValidateParams(t *testing.T) {
    // Valid cron
    c := getTestCron()

    // Valid parameters
    if c.ValidateParams() == false {
        t.Errorf("Expected ValidateParams() to return true, when passing valid parameters")
    }

    // Invalid parameters
    c = getTestCron()
    c.Account = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Account")
    }

    c = getTestCron()
    c.Cronname = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Cronname")
    }

    c = getTestCron()
    c.Email = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Email")
    }

    c = getTestCron()
    c.Frequency = 0
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Frequency")
    }

    c = getTestCron()
    c.Ipaddress = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad IP Address")
    }
    
    c = getTestCron()
    c.Lastruntime = -1
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Lastruntime")
    }

}

func TestCheckLength(t *testing.T) {
    // Valid cron
    c := getTestCron()

    // Valid parameters
    if err := c.CheckLength(); err != nil {
        t.Errorf("Expected ValidateParams() to return nil, when passing valid parameters, got:\n" + err.Error())
    }

    // Invalid parameters
    c.Account = ""
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad Account")
    }
    c.Cronname = ""
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad Cronname")
    }
    c.Email = ""
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad Email")
    }
    c.Frequency = 0
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad Frequency")
    }
    c.Ipaddress = ""
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad IP Address")
    }
    c.Lastruntime = -1
    if c.CheckLength() == nil {
        t.Errorf("Expected ValidateParams() to return an error, when using bad Lastruntime")
    }
}

func TeststringToInt(t *testing.T) {
    if stringToInt("w") != -1 {
        t.Errorf("Expected stringToInt() to return -1 when an invalid int was passed")
    }

    if stringToInt("1") != 1 {
        t.Errorf("Expected stringToInt() to return the int '1' when string \"1\" was passed")
    }
}
