package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RegistryOptions struct {
	AccessKey    string
	LocalAddress string
	BindAddress  string
	DefaultPort  string
}

func ResolveRegistryOptions(opts ...RegistryOptions) RegistryOptions {
	var options RegistryOptions

	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = RegistryOptions{}
	}

	var defaultPort string = options.DefaultPort
	if len(defaultPort) < 1 {
		defaultPort = "3000"
	}

	var localAddress string = options.LocalAddress

	if len(localAddress) < 1 {
		localAddress = "localhost:" + defaultPort
	}

	var bindAddress string = options.BindAddress

	if len(bindAddress) < 1 {
		hostname, err := os.Hostname()

		if err == nil {
			bindAddress = hostname
		}
	}

	return RegistryOptions{
		AccessKey:    options.AccessKey,
		LocalAddress: localAddress,
		BindAddress:  bindAddress,
		DefaultPort:  defaultPort,
	}
}

func (c client) GenerateRegistryPayload() ([]byte, error) {
	routes := map[string]interface{}{}

	for _, route := range c.Handlers {
		routes[string(route.Method)+"@"+route.Path] = map[string]string{
			"method": string(route.Method),
			"path":   route.Path,
		}
	}

	payload := map[string]interface{}{
		"device":       c.device,
		"name":         c.name,
		"title":        c.title,
		"description":  c.description,
		"version":      c.version,
		"mountPoint":   c.mountPoint,
		"mountVersion": c.mountVersion,
		"category":     c.category,
		"tags":         c.tags,
		"routes":       routes,
	}

	return json.Marshal(payload)
}

func (c client) Register(endpoint string, options RegistryOptions) {
	c.Logger.Info().Msgf("Registering client: %s", c.DisplayName())

	_, err := c.GenerateRegistryPayload()
	if err != nil {
		c.Logger.Panic().Str("error", err.Error()).Msg("Failed to generate registry payload")
	} else {
		// fmt.Printf("RESULT: %v\n", payload)
		c.Logger.Info().Msgf("Gateway mount path: %s", c.MountPath())

		if len(options.BindAddress) < 1 {
			c.Logger.Panic().Msg("Failed to resolve binding address")
		}

		c.Logger.Info().Msgf("Connected to registry gateway [%s]", options.BindAddress)

		requestURL := fmt.Sprintf("%s/register", endpoint)
		response, err := http.NewRequest(http.MethodGet, requestURL, nil)

		if err != nil {
			c.Logger.Panic().Str("error", err.Error()).Msg("client: could not create request")
		} else {
			fmt.Printf("RESPONSE: %v\n", response)
		}
	}
}
