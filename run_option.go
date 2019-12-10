package bobo

import (
	"net/http"

	"github.com/eure/bobo/command"
	"github.com/eure/bobo/engine"
	"github.com/eure/bobo/log"
)

// RunOption contains options.
type RunOption struct {
	Engine             engine.Engine
	CommandSet         *command.CommandSet
	Logger             log.Logger
	HTTPClient         *http.Client
	MaxRunningCommands int
	NoPanic            bool
}
