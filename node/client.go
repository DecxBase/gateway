package node

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobeam/stringy"
	"github.com/phuslu/log"
)

type ClientOptions struct {
	Device   string
	Logger   *log.Logger
	Timezone *time.Location
}

func (co ClientOptions) GetDevice() string {
	if len(co.Device) > 0 {
		return co.Device
	}

	device, err := os.Hostname()
	if err != nil {
		device = "localhost"
	}

	return device
}

func (co ClientOptions) GetTimezone() *time.Location {
	if co.Timezone != nil {
		return co.Timezone
	}

	return time.Local
}

type client struct {
	device       string
	name         string
	title        string
	description  string
	version      string
	mountPoint   string
	mountVersion string
	category     string
	tags         []string
	Logger       *log.Logger
	Handlers     []*clientRoute
	Router       *gin.Engine
}

func (c *client) SetDevice(device string) *client {
	c.device = device
	c.Logger.Context = log.NewContext(nil).Str("device", device).Value()
	return c
}

func (c *client) SetTitle(title string) *client {
	c.title = title
	return c
}

func (c *client) SetDescription(desc string) *client {
	c.description = desc
	return c
}

func (c *client) SetVersion(version string) *client {
	c.version = version
	return c
}

func (c *client) SetMount(path string) *client {
	if !ValidIdentifier(path) {
		c.Logger.Panic().Msgf("Invalid client mount point: %s", path)
	}

	c.mountPoint = path
	return c
}

func (c *client) SetMountVersion(version string) *client {
	if !ValidIdentifier(version) {
		c.Logger.Panic().Msgf("Invalid client mount version: %s", version)
	}

	c.mountVersion = version
	return c
}

func (c *client) SetCategory(cat string) *client {
	c.category = cat
	return c
}

func (c *client) SetTags(tags []string) *client {
	c.tags = tags
	return c
}

func (c client) DisplayName() string {
	if len(c.title) > 0 {
		return c.title
	}

	return stringy.New(c.name).Title()
}

func (c client) MountPath() string {
	var path = c.mountPoint

	if len(path) < 1 {
		path = c.name
	}

	if len(c.mountVersion) > 0 {
		path += "/" + c.mountVersion
	}

	return "/" + path
}

func (c client) Run(endpoint string, opts ...RegistryOptions) {
	var options = ResolveRegistryOptions(opts...)

	c.Register(endpoint, options)
	c.RunRouter(options.LocalAddress)
}

func (c client) String() string {
	return c.DisplayName() + "@" + c.version
}
