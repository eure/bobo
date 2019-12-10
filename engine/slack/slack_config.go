package slack

import (
	"fmt"
	"os"

	"github.com/eure/bobo/command"
	"github.com/eure/bobo/log"
)

type Config interface {
	GetLogger() log.Logger
	GetCommandSet() *command.CommandSet
	GetSlackToken() string
	GetMaxRunningCommands() int
}

func getToken(c Config) string {
	token := c.GetSlackToken()
	if token != "" {
		return token
	}

	for _, envName := range envSlackTokens {
		token = os.Getenv(envName)
		if token != "" {
			return token
		}
	}

	panic(fmt.Sprintf("[PANIC] Slack Token is empty. Use one of these envvar: %+v", envSlackTokens))
}

// envvar list for Slack Token.
var envSlackTokens = []string{
	"SLACK_RTM_TOKEN",
	"SLACK_BOT_TOKEN",
	"SLACK_TOKEN",
}
