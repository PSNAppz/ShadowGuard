package receiver

import (
	"fmt"
)

// NotificationReceiver is an interface that represents the receivers of the notification.
type NotificationReceiver interface {
	Type() string
	Notify(message string) error
}

func CreateReceivers(pluginSettings map[string]interface{}) ([]NotificationReceiver, error) {
	receiverSettings, err := ParseReceiverSettings(pluginSettings)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("notifcation setting must specify a type, found %v", notificationType)
	}

	var receiver NotificationReceiver
	var err error

	switch notificationType {
	case "slack":
		receiver, err = NewSlackReceiver(setting)
	case "file":
		receiver, err = NewFileReceiver(setting)
	}

	return receiver, err
}

func ParseReceiverSettings(settings map[string]interface{}) ([]map[string]interface{}, error) {
	notificationList, ok := settings["notify"]
	if !ok {
		return nil, nil
	}

	notificationInterfaceList, ok := notificationList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("notification list is incorrectly configured, found %+v", notificationList)
	}

	notificationSettings := []map[string]interface{}{}
	for _, notificationInterface := range notificationInterfaceList {
		notificationSetting, ok := notificationInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("notification settting is incorrectly configured, found %+v", notificationSetting)
		}
		notificationSettings = append(notificationSettings, notificationSetting)
	}

	return notificationSettings, nil
}
