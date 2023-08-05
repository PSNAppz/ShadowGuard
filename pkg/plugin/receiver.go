package plugin

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

func (s SlackReceiver) ReceiveNotification(message string) error {
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
		return nil, fmt.Errorf("channel ID must be a string")
	}

	token, ok := settings["token"]
	if !ok {
		return nil, fmt.Errorf("API token required to send Slack notifications")
	}

	tokenStr, ok := token.(string)
	if !ok {
		return nil, fmt.Errorf("API token must be a string")
	}

	fmt.Println("token", tokenStr)
	fmt.Println("channel", channelIDStr)
	return &SlackReceiver{api: slack.New(tokenStr, slack.OptionDebug(true)), channelID: channelIDStr}, nil
}

func ParseReceivers(settings []map[string]interface{}) ([]NotificationReceiver, error) {
	receivers := []NotificationReceiver{}
	for _, setting := range settings {
		notificationType, ok := setting["type"]
		if !ok {
			return nil, fmt.Errorf("notifcation setting must specify a type")
		}

		switch notificationType {
		case "slack":
			receiver, err := NewSlackReceiver(setting)
			if err != nil {
				return nil, err
			}
			receivers = append(receivers, receiver)
		}
	}
	return receivers, nil
}

// NotificationReceiver is an interface that represents the receivers of the notification.
type NotificationReceiver interface {
	Type() string
	ReceiveNotification(message string) error
}

func SendNotification(message string, receivers ...NotificationReceiver) {
	for _, receiver := range receivers {
		err := receiver.ReceiveNotification(message)
		if err != nil {
			// Handle error if needed
			fmt.Printf("Error sending notification: %v\n", err)
		}
	}
}
