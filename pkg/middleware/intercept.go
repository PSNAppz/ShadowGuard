package middleware

import (
	"shadowguard/pkg/config"
	"shadowguard/pkg/database"
	"shadowguard/pkg/plugin"
	_ "shadowguard/plugins" // Import the plugins package to register the plugins

	"io"
	"log"
	"net/http"
)

// Intercept performs intercept operation, contacts internal server and returns response to client
func Intercept(client *http.Client, method, url string, pluginConfigs []config.PluginConfig, db database.DB) http.HandlerFunc {
	// Convert PluginConfigs to Plugins
	plugins, err := createPlugins(pluginConfigs, db)
	if err != nil {
		panic(err)
	}
	// Separate plugins into active and passive
	activePlugins := []plugin.Plugin{}
	passivePlugins := []plugin.Plugin{}
	for _, p := range plugins {
		if p.IsActiveMode() {
			activePlugins = append(activePlugins, p)
		} else {
			passivePlugins = append(passivePlugins, p)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() // close initial request

		// Execute passive plugins in separate goroutines
		for _, p := range passivePlugins {
			go p.Handle(r)
		}

		// Execute active plugins
		for _, p := range activePlugins {
			err := p.Handle(r)
			if err != nil {
				log.Printf("Request blocked by plugin %s. Error %s", p.Type(), err.Error())
				// If an active plugin returns an error, respond with an error message and status code
				http.Error(w, "Request blocked by plugin: "+err.Error(), http.StatusForbidden)
				return
			}
		}

		// route request to internal server
		req, err := http.NewRequest(method, url, r.Body)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		w.Write(respBody)
	}
}

func createPlugins(pluginConfigs []config.PluginConfig, db database.DB) ([]plugin.Plugin, error) {
	plugins := make([]plugin.Plugin, len(pluginConfigs))
	for i, pc := range pluginConfigs {
		p, err := plugin.CreatePlugin(pc.Type, pc.Settings, db)
		if err != nil {
			return nil, err
		}
		plugins[i] = p
	}
	return plugins, nil
}
