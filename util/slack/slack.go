package slack

import (
    "bytes"
    "strconv"
    "encoding/json"
    "net/http"
    "errors"

    "gocron/util/log"
)

// Slack struct defines a slack object, and has access to
// the Message() funtion
type Slack struct {
    HookURL string
    Post    APIPost
}

// APIPost struct defines a slack post
type APIPost struct {
    Channel string `json:"channel"`
    Text    string `json:"text"`
}

// Message sends a slack message
func (s Slack) Message() error {
    if err := s.validate(); err != nil {
        return err
    }

    payload, err := json.Marshal(s.Post)
    if err != nil {
        return err
    }
    log.Message(string(payload))
    return s.sendPayload(payload)
}

func (s Slack) sendPayload(p []byte) error {
    req, err := http.NewRequest("POST", s.HookURL, bytes.NewBuffer(p))
    if err != nil {
        return err
    }
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Slack returned status: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func (s Slack) validate() error {
    if len(s.Post.Text) == 0 {
        return errors.New("slack message is empty")
    }
    return nil
}
