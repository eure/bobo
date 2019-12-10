package bobo

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/eure/bobo/engine"
	"github.com/eure/bobo/errorcode"
	"github.com/eure/bobo/log"
)

// Bot is core struct for a bot.
type Bot struct {
	logger log.Logger
	engine engine.Engine

	// flags
	status bool

	signal chan os.Signal
}

// NewBot returns initialized Bot.
func NewBot() (*Bot, error) {
	return NewBotWithConfig(Config{})
}

// NewBot returns initialized Bot from Config.
func NewBotWithConfig(conf Config) (*Bot, error) {
	if conf.CommandSet == nil {
		return nil, errors.New("You must set CommandSet for bot")
	}

	e := conf.Engine
	if e == nil {
		return nil, errors.New("You must set an engine for bot platform")
	}
	err := e.Init(conf)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		engine: e,
		logger: conf.GetLogger(),
		status: true,
		signal: make(chan os.Signal, 1),
	}
	go bot.catchSignal()

	return bot, nil
}

// Run starts to run engine.
func (b *Bot) Run() (errCode int) {
	return b.engine.Run()
}

// SetStatus sets active status.
// if it's false, bot does not react to any message from a user.
func (b *Bot) SetStatus(f bool) {
	b.status = f
}

// LogInfo logs info level log.
func (b Bot) LogInfo(typ, msg string, v ...interface{}) {
	b.logger.Infof(typ, msg, v...)
}

// catchSignal handles OS signals.
func (b *Bot) catchSignal() {
	signal.Notify(b.signal,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		s := <-b.signal
		switch s {
		case syscall.SIGHUP:
			b.LogInfo("catchSignal", "syscall.SIGHUP")
			return
		case syscall.SIGINT:
			b.LogInfo("catchSignal", "syscall.SIGINT")
			b.exit()
			return
		case syscall.SIGTERM:
			b.LogInfo("catchSignal", "syscall.SIGTERM")
			b.exit()
			return
		case syscall.SIGQUIT:
			b.LogInfo("catchSignal", "syscall.SIGQUIT")
			b.exit()
			return
		}
	}
}

// func (b *Bot) reload() {
// 	b.LogInfo("reload", "reloading...")
// 	b.engine.Reload()
// 	signal.Stop(b.signal)
// }

func (b *Bot) exit() {
	b.LogInfo("exit", "closing...")
	signal.Stop(b.signal)
	b.engine.Close(errorcode.Exit)
}
