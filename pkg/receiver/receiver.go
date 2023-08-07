package receiver

import (
	"fmt"

	"github.com/slack-go/slack"
)

// NotificationReceiver is an interface that represents the receivers of the notification.
type NotificationReceiver interface {
	Type() string
	Notify(message string) error
}

func CreateReceivers(pluginSettings map[string]interface{}) ([]NotificationReceiver, error) {
	receiverSettings := ParseReceiverSettings(pluginSettings)

	receivers := []NotificationReceiver{}
	for _, setting := range receiverSettings {
		receiver, err := NewReciever(setting)
		if err != nil {
			return nil, err

		}
		receivers = append(receivers, receiver)
	}
	return receivers, nil
}

func NewReciever(setting map[string]interface{}) (NotificationReceiver, error) {
	notificationType, ok := setting["type"]
	if !ok {
		return nil, fmt.Errorf("notifcation setting must specify a type")
	}

	var receiver NotificationReceiver
	var err error

	switch notificationType {
	case "slack":
		receiver, err = NewSlackReceiver(setting)
		if err != nil {
			return nil, err
		}
		return receiver, nil
	}

	return receiver, err
}

func ParseReceiverSettings(settings map[string]interface{}) []map[string]interface{} {
	if _, ok := settings["notify"]; !ok {
		return nil
	}

	notificationInterfaceList, ok := settings["notify"].([]interface{})
	if !ok {
		panic("")
		//panic(fmt.Errorf("notification list incorrectly configured, found %+v", notificationSettingsInterface))
	}

	notificationSettings := []map[string]interface{}{}
	for _, notificationInterface := range notificationInterfaceList {
		notificationSetting, ok := notificationInterface.(map[string]interface{})
		if !ok {
			panic("")
			//panic(fmt.Errorf("notification list incorrectly configured, found %+v", notificationSettingsInterface))
		}
		notificationSettings = append(notificationSettings, notificationSetting)
	}

	return notificationSettings
}

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

	return &SlackReceiver{api: slack.New(tokenStr, slack.OptionDebug(true)), channelID: channelIDStr}, nil
}
