package main

import (
	"github.com/eure/bobo"
	"github.com/eure/bobo/command"
	"github.com/eure/bobo/engine/slack"
	"github.com/eure/bobo/log"
)

// Entry Point
func main() {
	bobo.Run(bobo.RunOption{
		Engine: &slack.SlackEngine{},
		Logger: &log.StdLogger{
			IsDebug: bobo.IsDebug(),
		},
		CommandSet: command.NewCommandSet(
			command.PingCommand,
			command.ParrotCommand,
			command.GoodMorningCommand,
			command.ReactEmojiCommand,
			command.ReloadCommand,
			command.HelpCommand,
		),
	})
}
