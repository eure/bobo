package bobo

import (
	"net/http"

	"github.com/eure/bobo/command"
	"github.com/eure/bobo/engine"
	"github.com/eure/bobo/log"
)

// Config is a config struct for bot.
type Config struct {
	engine.Engine
	*command.CommandSet
	Logger     log.Logger
	HTTPClient *http.Client

	SlackToken         string
	MaxRunningCommands int
}

func (c Config) GetCommandSet() *command.CommandSet {
	return c.CommandSet
}

func (c Config) GetSlackToken() string {
	return c.SlackToken
}

func (c Config) GetLogger() log.Logger {
	if c.Logger != nil {
		return c.Logger
	}

	return log.DefaultLogger
}

func (c Config) GetMaxRunningCommands() int {
	if c.MaxRunningCommands != 0 {
		return c.MaxRunningCommands
	}

	const defaultMaxRunning = 5
	return defaultMaxRunning
}
