package registry

import (
	"context"
	"errors"

	"net"
	"net/url"
	"os"

	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
)

type LBStrategy int

const (
	RoundRobin LBStrategy = iota
	LeastConnected
)

var Logger log.Logger

func InitLogger() log.Logger {
	Logger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
			// Writer:         os.Stdout,
			// Formatter: log.LogfmtFormatter{"ts"}.Formatter,
		},
	}
	return Logger
}

func GetLBStrategy(strategy string) LBStrategy {
	switch strategy {
	case "least-connection":
		return LeastConnected
	default:
		return RoundRobin
	}
}

type Config struct {
	Port            int      `yaml:"lb_port"`
	MaxAttemptLimit int      `yaml:"max_attempt_limit"`
	Backends        []string `yaml:"backends"`
	Strategy        string   `yaml:"strategy"`
}

const MAX_LB_ATTEMPTS int = 3

func GetLBConfig() (*Config, error) {
	var config Config
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	if len(config.Backends) == 0 {
		return nil, errors.New("backend hosts expected, none provided")
	}

	if config.Port == 0 {
		return nil, errors.New("load balancer port not found")
	}

	return &config, nil
}

func IsBackendAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		Logger.Debug().Str("error", err.Error()).Msg("Site unreachable")
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}
