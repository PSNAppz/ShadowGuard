package publisher

import (
	"fmt"
)

// Publisher is an interface that publishes messages to a data source
type Publisher interface {
	Type() string
	Publish(message string) error
}

func CreatePublishers(pluginSettings map[string]interface{}) ([]Publisher, error) {
	publisherSettings, err := parsePublisherSettings(pluginSettings)
	if err != nil {
		return nil, err
	}

	publishers := []Publisher{}
	for _, setting := range publisherSettings {
		publisher, err := NewPublisher(setting)
		if err != nil {
			return nil, err

		}
		publishers = append(publishers, publisher)
	}
	return publishers, nil
}

func NewPublisher(setting map[string]interface{}) (Publisher, error) {
	publisherType, ok := setting["type"]
	if !ok {
		return nil, fmt.Errorf("publish setting must specify a type, found %v", publisherType)
	}

	var publisher Publisher
	var err error

	switch publisherType {
	case "slack":
		publisher, err = NewSlackPublisher(setting)
	case "file":
		publisher, err = NewFilePublisher(setting)
	}

	return publisher, err
}

func parsePublisherSettings(settings map[string]interface{}) ([]map[string]interface{}, error) {
	publisherList, ok := settings["publishers"]
	if !ok {
		return nil, nil
	}

	publisherInterfaceList, ok := publisherList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("publisher list is incorrectly configured, found %+v", publisherList)
	}

	publisherSettings := []map[string]interface{}{}
	for _, publisherInterface := range publisherInterfaceList {
		publisherSetting, ok := publisherInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("publisher settting is incorrectly configured, found %+v", publisherSetting)
		}
		publisherSettings = append(publisherSettings, publisherSetting)
	}

	return publisherSettings, nil
}
