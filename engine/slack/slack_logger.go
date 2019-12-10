package slack

import (
	"github.com/eure/bobo/log"
)

type SlackLogger struct {
	logger log.Logger
}

func (l SlackLogger) Output(_ int, s string) error {
	l.logger.Infof("SlackLogger", s)
	return nil
}
