package bobo

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jpillora/overseer"

	"github.com/eure/bobo/errorcode"
)

// Run is entry-point of bot daemon.
func Run(opt RunOption) {
	switch {
	case opt.UseUpgrade:
		runWithUpgrading(opt)
	default:
		run(opt)
	}
}

// run is entry-point of bot daemon.
func run(opt RunOption) {
	logger := opt.Logger
	logger.Infof("bobo", "pid=[%d]", os.Getpid())
	for {
		errCode := runBot(opt)
		logger.Errorf("bobo", "errCode=[%d]", errCode)
		if errCode != errorcode.None {
			os.Exit(errCode)
		}
		time.Sleep(time.Second * 5)
		logger.Infof("bobo", "run again...")
	}
}

// runWithUpgrading is entry-point of bot daemon with self-upgrading binary.
// ref: https://github.com/jpillora/overseer
func runWithUpgrading(opt RunOption) {
	overseer.Run(overseer.Config{
		Program: func(state overseer.State) {
			run(opt)
		},
		Fetcher: opt.UpgradeFetcher,
		Debug:   opt.UpgradeDebug,
	})
}

func runBot(opt RunOption) (errCode int) {
	bot, err := NewBotWithConfig(Config{
		Engine:             opt.Engine,
		CommandSet:         opt.CommandSet,
		Logger:             opt.Logger,
		HTTPClient:         opt.HTTPClient,
		MaxRunningCommands: opt.MaxRunningCommands,
	})
	if err != nil {
		panic(err)
	}

	if opt.NoPanic {
		defer func() {
			if err := recover(); err != nil {
				bot.LogInfo("PANIC", fmt.Sprintf("%v", err))
				bot.exit()
				errCode = errorcode.None
			}
		}()
	}

	return bot.Run()
}

// IsDebug checks debug flag is true or not.
// debug flag is set via envvar "BOBO_DEBUG".
func IsDebug() bool {
	b, _ := strconv.ParseBool(os.Getenv("BOBO_DEBUG"))
	return b
}
