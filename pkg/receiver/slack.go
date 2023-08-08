package receiver

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackReceiver struct {
	api       *slack.Client
	channelID string
}

func (s SlackReceiver) Type() string {
	return "slack"
}

func (s SlackReceiver) Notify(message string) error {
	_, _, err := s.api.PostMessage(s.channelID, slack.MsgOptionText(message, false))
	return err
}

func NewSlackReceiver(settings map[string]interface{}) (*SlackReceiver, error) {
	channelID, ok := settings["channelID"]
	if !ok {
		return nil, fmt.Errorf("channel ID is required to send Slack notification")
	}

	channelIDStr, ok := channelID.(string)
	if !ok {
		return nil, fmt.Errorf("channel ID must be a string, found: %+v", channelID)
	}

	token, ok := settings["token"]
	if !ok {
		return nil, fmt.Errorf("API token required to send Slack notifications")
	}

	tokenStr, ok := token.(string)
	if !ok {
		return nil, fmt.Errorf("API token must be a string, found: %+v", tokenStr)
	}

	return &SlackReceiver{api: slack.New(tokenStr, slack.OptionDebug(true)), channelID: channelIDStr}, nil
}
