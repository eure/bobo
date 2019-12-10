package command

import (
	"fmt"
	"regexp"
)

// GoodMorningCommand is an example command.
// This command reply greeting to a thread.
var GoodMorningCommand = BasicCommandTemplate{
	Help:   "reply Good morning in thread",
	Regexp: regexp.MustCompile("(?i)Good Morning"),
	GenerateFn: func(d CommandData) Command {
		c := Command{}
		text := fmt.Sprintf("Good Morning <@%s>", d.SenderID)
		task := NewReplyThreadEngineTask(d.Engine, d.Channel, text, d.ThreadTimestamp)
		c.Add(task)
		return c
	},
}
