package command

import (
	"fmt"
	"strings"
)

// HelpCommand is an example command.
// This command shows all of the commands help text.
var HelpCommand = BasicCommandTemplate{
	Help:           "show command list and usages",
	MentionCommand: "help",
	GenerateFn: func(d CommandData) Command {
		c := Command{}

		mentionMap := d.GetMentionMap()
		mentions := make([]string, 0, len(mentionMap))
		for key, c := range mentionMap {
			if !c.HasHelp() {
				continue
			}
			str := fmt.Sprintf("- %s\t\t%s", key, c.GetHelp())
			mentions = append(mentions, str)
		}

		regexpList := d.GetRegexpList()
		freewords := make([]string, len(regexpList))
		for i, c := range regexpList {
			if !c.HasHelp() {
				continue
			}
			freewords[i] = fmt.Sprintf("- %s\t\t%s", c.GetRegexp().String(), c.GetHelp())
		}

		text := fmt.Sprintf("```Command:\n%s\n\nFreeword:\n%s```", strings.Join(mentions, "\n"), strings.Join(freewords, "\n"))
		task := NewReplyEngineTask(d.Engine, d.Channel, text)
		c.Add(task)
		return c
	},
}
