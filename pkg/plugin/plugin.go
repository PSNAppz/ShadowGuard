package plugin

import (
	"fmt"
	"net/http"
	"shadowguard/pkg/database"
)

// Plugin is an interface that any plugin should implement.
type Plugin interface {
	Type() string
	IsActiveMode() bool
	Notify(message string)
	Handle(r *http.Request) error
}

// PluginFactory is a function that creates a Plugin.
type PluginFactory func(settings map[string]interface{}, db *database.Database) Plugin

// pluginRegistry holds a map of plugin types to factory functions.
var pluginRegistry = make(map[string]PluginFactory)

// RegisterPlugin registers a new type of plugin with a factory function.
func RegisterPlugin(typeStr string, factory PluginFactory) {
	pluginRegistry[typeStr] = factory
}

// CreatePlugin creates a new Plugin of the given type with the given settings.
func CreatePlugin(typeStr string, settings map[string]interface{}, db *database.Database) (Plugin, error) {
	factory, exists := pluginRegistry[typeStr]
	if !exists {
		return nil, fmt.Errorf("plugin type %s not found in registry", typeStr)
	}
	return factory(settings, db), nil
}
