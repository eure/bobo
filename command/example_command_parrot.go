package command

import (
	"fmt"
)

// ParrotCommand is an example command.
// This command puts :parrot: emoji both side.
var ParrotCommand = BasicCommandTemplate{
	Help:           "reply text between :parrot: Emoji",
	MentionCommand: "parrot",
	GenerateFn: func(d CommandData) Command {
		c := Command{}
		text := fmt.Sprintf(":parrot: %s :parrot:", d.TextOther)
		task := NewReplyEngineTask(d.Engine, d.Channel, text)
		c.Add(task)
		return c
	},
}
