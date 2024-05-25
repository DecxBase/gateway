package node

import (
	"fmt"
	"regexp"

	"github.com/phuslu/log"
)

func ValidIdentifier(txt string) bool {
	matched, _ := regexp.MatchString("^([a-zA-Z0-9_]+)$", txt)

	return matched
}

func Client(name string, options *ClientOptions) *client {
	var logger *log.Logger

	if !ValidIdentifier(name) {
		panic(fmt.Sprintf("Invalid client name: %s", name))
	}

	if options.Logger != nil {
		logger = options.Logger
	} else {
		logger = &log.Logger{
			Caller:       1,
			TimeFormat:   "15:04:05",
			TimeLocation: options.GetTimezone(),
			Context:      log.NewContext(nil).Str("device", options.GetDevice()).Value(),
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	return &client{
		name:     name,
		device:   options.GetDevice(),
		Logger:   logger,
		Handlers: []*clientRoute{},
	}
}

func DefaultClient(name string) *client {
	return Client(name, &ClientOptions{})
}
