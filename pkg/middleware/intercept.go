package middleware

import (
	"AegisGuard/pkg/config"
	"AegisGuard/pkg/plugin"
	_ "AegisGuard/plugins" // Import the plugins package to register the plugins
	"log"

	"io"
	"net/http"
)

// Intercept performs intercept operation, contacts internal server and returns response to client
func Intercept(client *http.Client, method, url string, pluginConfigs []config.PluginConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Convert PluginConfigs to Plugins
		plugins, err := createPlugins(pluginConfigs)
		if err != nil {
			http.Error(w, "Error creating plugins: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute plugins
		for _, p := range plugins {
			if p.GetMode() == plugin.Passive {
				// Execute passive plugins in separate goroutines
				go p.Handle(r)
			} else {
				// Execute active plugins
				err := p.Handle(r)
				if err != nil {
					// If an active plugin returns an error, log the error message and status code
					log.Printf("Request blocked by plugin %s. Error %s", p.GetType(), err.Error())
				}
			}
		}

		defer r.Body.Close()
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

func createPlugins(pluginConfigs []config.PluginConfig) ([]plugin.Plugin, error) {
	plugins := make([]plugin.Plugin, len(pluginConfigs))
	for i, pc := range pluginConfigs {
		p, err := plugin.CreatePlugin(pc.Type, pc.Settings)
		if err != nil {
			return nil, err
		}
		plugins[i] = p
	}
	return plugins, nil
}
