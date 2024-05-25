package node

import (
	"fmt"
	"time"

	"github.com/phuslu/log"
)

type ClientOptions struct {
	Timezone *time.Location
	Logger   *log.Logger
}

func (co ClientOptions) GetTimezone() *time.Location {
	if co.Timezone != nil {
		return co.Timezone
	}

	return time.Local
}

type client struct {
	Name         string
	Title        string
	Description  string
	Version      string
	MountPoint   string
	MountVersion string
	Logger       *log.Logger
}

func (c *client) SetTitle(title string) *client {
	c.Title = title
	return c
}

func (c *client) SetDescription(desc string) *client {
	c.Description = desc
	return c
}

func (c *client) SetVersion(version string) *client {
	c.Version = version
	return c
}

func (c *client) SetMount(path string) *client {
	c.MountPoint = path
	return c
}

func (c *client) SetMountVersion(version string) *client {
	c.MountVersion = version
	return c
}

func (c client) Initiate(endpoint string, accessKey string) {
	fmt.Printf("TIME TO INITIATE CLIENT: %v\n", c)
}

func Client(name string, options *ClientOptions) *client {
	var logger *log.Logger

	if options.Logger != nil {
		logger = options.Logger
	} else {
		logger = &log.Logger{
			Caller:       1,
			TimeFormat:   "15:04:05",
			TimeLocation: options.GetTimezone(),
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	return &client{
		Name:   name,
		Logger: logger,
	}
}

func DefaultClient(name string) *client {
	return Client(name, &ClientOptions{})
}
