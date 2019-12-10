package command

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

// ReactEmojiCommand is an example command.
// This command add emoji reaction 20% chance.
var ReactEmojiCommand = BasicCommandTemplate{
	Help:   "add reaction to message",
	Regexp: regexp.MustCompile("^"),
	GenerateFn: func(d CommandData) Command {
		c := Command{}
		if !isRandValid(20) {
			return c
		}

		emoji, err := d.Engine.GetEmojiByRandom()
		if err != nil {
			task := NewReplyEngineTask(d.Engine, d.Channel, fmt.Sprintf("Error on Slack!\nErr: `%s`", err.Error()))
			c.Add(task)
			return c
		}

		task := NewReactionEmojiEngineTask(d.Engine, d.Channel, emoji, d.ThreadTimestamp)
		c.Add(task)
		return c
	},
}

func isRandValid(percent int) bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return percent > rand.Intn(100)
}
