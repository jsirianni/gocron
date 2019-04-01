package slack
import (
    "testing"
)

func TestMessage(t *testing.T) {
    var s Slack

    // Test empty struct
    if err := s.Message(); err == nil {
        t.Errorf("Expected Message() to return an error due to an empty Slack struct")
    }

    // Test URL
    s.Post.Channel = "test"
    s.Post.Text    = "test"

    s.HookURL = "http://badurl.com"
    if err := s.Message(); err == nil {
        t.Errorf("Expected Message() to return an error when using a bad url")
    }

    s.HookURL = "https://httpstat.us/200"
    if err := s.Message(); err != nil {
        t.Errorf("Expected Message() to return a nil error, instead got: " + err.Error())
    }

    // Test Channel
    s.Post.Text = ""
    if err := s.Message(); err == nil {
        t.Errorf("Expected Message() to return an error when using an empty Slack.Post.Text")
    }
}
