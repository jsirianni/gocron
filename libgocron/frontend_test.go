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
    c.Account = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Account")
    }
    c.Cronname = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Cronname")
    }
    c.Email = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Email")
    }
    c.Frequency = 0
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Frequency")
    }
    c.Ipaddress = ""
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad IP Address")
    }
    c.Lastruntime = -1
    if c.ValidateParams() == true {
        t.Errorf("Expected ValidateParams() to return false, when using bad Lastruntime")
    }

}
